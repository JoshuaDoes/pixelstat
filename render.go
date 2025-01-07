package main

import (
  "github.com/JoshuaDoes/crunchio"
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

  hz         time.Duration
  pollCnt    int64
}

func NewRenderer(hz int) *Renderer {
  r := new(Renderer)
  r.NodeBase = NewNodeBase()
  r.hz = hertz(time.Duration(hz))
  terminal = ncurses.Init()
  return r
}

func (r *Renderer) Close() {
  ncurses.EndWin()
  terminal = nil
}

func (r *Renderer) Value() *crunchio.Buffer {
  r.Lock()
  defer r.Unlock()
  if terminal == nil {
    return nil
  }
  
  nodes := r.Nodes()
  pollCnt := int64(0)
  for i := 0; i < len(nodes); i++ {
    n := nodes[i]
    if n.Name() == "renderer" {
      continue
    }
    pollCnt += n.GetTracker().Value(n.Name()).Polls()
  }
  if pollCnt <= r.pollCnt {
    return nil
  }
  r.pollCnt = pollCnt

  nameLen := 0
  unitLen := 0
  for i := 0; i < len(nodes); i++ {
    n := nodes[i]
    if n.Name() == "renderer" {
      continue
    }
    if nl := len(n.Name()); nl > nameLen {
      nameLen = nl
    }
    if ul := len(n.Unit()); ul > unitLen {
      unitLen = ul
    }
  }

  nameF := "%" + fmt.Sprintf("%d", nameLen) + "s"
  live := "Live:\n"
  avg := "Average:\n"
  for i := 0; i < len(nodes); i++ {
    n := nodes[i]
    if n.Name() == "renderer" {
      continue
    }
    u := n.Unit()
    tv := n.GetTracker().Value(n.Name())
    v := tv.Value()
    if v.ByteCapacity() == 0 {
      continue
    }

    vs := ""
    as := ""
    vt := n.ValueType()
    switch vt {
    case "raw":
      length := n.ValueLen()
      if bc := v.ByteCapacity(); bc < length {
        length = bc
      }
      vs = fmt.Sprintf("% X", v.ReadBytes(0, length)) + u
    case "str":
      length := n.ValueLen()
      if bc := v.ByteCapacity(); bc < length {
        length = bc
      }
      vs = string(v.ReadBytes(0, length)) + u
    case "i32":
      length := n.ValueLen()
      if bc := v.ByteCapacity(); (bc / 4) < length {
        length = (bc / 4)
      }
      for i := int64(0); i < length; i++ {
        if i > 0 {
          vs += " "
          as += " "
        }
        vs += i32tostr(v.ReadI32LENext(1)[0]) + u
        as += f64tostr(tv.Average(i), 0) + u
      }
    case "i64":
      length := n.ValueLen()
      if bc := v.ByteCapacity(); (bc / 8) < length {
        length = (bc / 8)
      }
      for i := int64(0); i < length; i++ {
        if i > 0 {
          vs += " "
          as += " "
        }
        vs += i64tostr(v.ReadI64LENext(1)[0]) + u
        as += f64tostr(tv.Average(i), 0) + u
      }
    case "f32":
      length := n.ValueLen()
      if bc := v.ByteCapacity(); (bc / 4) < length {
        length = (bc / 4)
      }
      for i := int64(0); i < length; i++ {
        if i > 0 {
          vs += " "
          as += " "
        }
        vs += f32tostr(v.ReadF32LENext(1)[0], 2) + u
        as += f64tostr(tv.Average(i), 2) + u
      }
    case "f64":
      length := n.ValueLen()
      if bc := v.ByteCapacity(); (bc / 8) < length {
        length = (bc / 8)
      }
      for i := int64(0); i < length; i++ {
        if i > 0 {
          vs += " "
          as += " "
        }
        vs += f64tostr(v.ReadF64LENext(1)[0], 2) + u
        as += f64tostr(tv.Average(i), 2) + u
      }
    default:
      perr(fmt.Errorf("renderer: invalid type: %s", vt))
    }

    live += fmt.Sprintf(nameF + ": %s\n", n.Name(), vs)
    if as != "" {
      avg += fmt.Sprintf(nameF + ": %s (%.2fpps, %d)\n", n.Name(), as, tv.PollsPerSecond(), tv.Polls())
    }
  }

  if terminal == nil {
    return nil
  }
  terminal.Erase()
  terminal.Printf("FPS: %.2f\n\n", r.GetTracker().Value(r.Name()).PollsPerSecond())
  terminal.Printf("%s\n", live)
  terminal.Printf("%s\n", avg)
  terminal.Refresh()

  return crunchio.NewBuffer(make([]byte, 1))
}

func (r *Renderer) Name() string {
  return "renderer"
}

func (r *Renderer) Rate() time.Duration {
  return r.hz
}

func (r *Renderer) ValueLen() int64 {
  return 1
}
