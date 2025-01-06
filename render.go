package main

import (
  "seehuhn.de/go/ncurses"

  "fmt"
  "sync"
  "time"
)

var (
  terminal *ncurses.Window
)

type Renderer struct {
  sync.Mutex
  *NodeBase
}

func NewRenderer() *Renderer {
  r := new(Renderer)
  r.NodeBase = NewNodeBase()
  return r
}

func (r *Renderer) RefreshTerminal(hz time.Duration) chan bool {
  terminal = ncurses.Init()
  stopper := loop(hertz(hz), r.terminal)
  go func(stopper chan bool, r *Renderer) {
    <- stopper
    r.Lock()
    defer r.Unlock()
    ncurses.EndWin()
    terminal = nil
  }(stopper, r)
  return stopper
}

func (r *Renderer) terminal() {
  if terminal == nil {
    return
  }

  nodes := r.Nodes()
  str := ""
  nameLen := 0
  unitLen := 0
  valueLen := 0

  for i := 0; i < len(nodes); i++ {
    n := nodes[i]
    if nl := len(n.Name()); nl > nameLen {
      nameLen = nl
    }
    if ul := len(n.Unit()); ul > unitLen {
      unitLen = ul
    }
    if vl := n.ValueLen(); vl > valueLen {
      valueLen = vl
    }
  }
  nameF := "%" + fmt.Sprintf("%d", nameLen) + "s"
  valueF := "%" + fmt.Sprintf("%d", valueLen) + "s"

  for i := 0; i < len(nodes); i++ {
    n := nodes[i]
    v := n.GetTracker().Value(n.Name())
    if v.ByteCapacity() == 0 {
      continue
    }

    u := n.Unit()
    ul := len(u)
    for i := 0; i < (unitLen - ul); i++ {
      u += " "
    }

    str += fmt.Sprintf(nameF + ": " + valueF + "%s\n", n.Name(), v, u)
  }

  r.Lock()
  defer r.Unlock()
  if terminal == nil {
    return
  }
  terminal.Erase()
  terminal.Printf("%s", str)
  terminal.Refresh()
}

func (r *Renderer) Name() string {
  return "renderer"
}
