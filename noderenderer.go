package main

import (
  "github.com/JoshuaDoes/crunchio"
  "seehuhn.de/go/ncurses"

  "fmt"
  "math"
  "sync"
  "time"
)

var (
  terminal *ncurses.Window
)

type NodeRenderer struct {
  sync.Mutex
  *NodeBase

  hz         time.Duration
  vrr        bool
  null       *crunchio.Buffer

  nameLen    int
  nameFmt    string
  unitLen    int

  polls      map[string]*NodeRender
}

func NewNodeRenderer(hz int, vrr bool) *NodeRenderer {
  n := new(NodeRenderer)
  n.NodeBase = NewNodeBase()
  n.hz = hertz(time.Duration(hz))
  n.vrr = vrr
  n.null = crunchio.NewBuffer(make([]byte, 1))
  n.polls = make(map[string]*NodeRender)
  terminal = ncurses.Init()
  return n
}

func (n *NodeRenderer) Poll(node string) bool {
  t := n.Node(node)
  if t == nil {
    return false
  }
  name := t.Name()
  if _, exists := n.polls[name]; !exists {
    r := new(NodeRender)
    r.p = n
    r.node = t
    n.polls[name] = r
  }
  return n.polls[name].Poll()
}

func (n *NodeRenderer) GetRender(node string) *NodeRender {
  if render, exists := n.polls[node]; exists {
    return render
  }
  return nil
}

func (n *NodeRenderer) SetTracker(t *Tracker) {
  nodes := t.Nodes()

  nameLen := 0
  unitLen := 0
  for i := 0; i < len(nodes); i++ {
    node := nodes[i]
    if node.ValueType() == "" {
      continue
    }
    if nl := len(node.Name()); nl > nameLen {
      nameLen = nl
    }
    if ul := len(node.Unit()); ul > unitLen {
      unitLen = ul
    }
  }

  n.nameLen = nameLen
  n.nameFmt = "%" + fmt.Sprintf("%d", nameLen) + "s"
  n.unitLen = unitLen
  n.Tracker = t
}

func (n *NodeRenderer) Close() error {
  ncurses.EndWin()
  terminal = nil
  return nil
}

func (n *NodeRenderer) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()
  if terminal == nil {
    return nil
  }

  nodes := n.Nodes()

  newFrame := true
  if n.vrr {
    newFrame = false
  }
  for i := 0; i < len(nodes); i++ {
    t := nodes[i]
    if t.ValueType() == "" {
      continue
    }
    if n.Poll(t.Name()) {
      newFrame = true
    }
  }
  if !newFrame {
    return nil
  }

  now := ""
  average := ""
  for i := 0; i < len(nodes); i++ {
    if render := n.GetRender(nodes[i].Name()); render != nil {
      if val := render.val; val != "" {
        now += val + "\n"
      }
      if avg := render.avg; avg != "" {
        average += avg + "\n"
      }
    }
  }

  if terminal == nil {
    return nil
  }
  terminal.Erase()
  terminal.Printf("FPS: %.0f\n\n", math.Round(n.GetTracker().Value(n.Name()).PollsPerSecond()))
  terminal.Printf("Now:\n%s\n", now)
  terminal.Printf("Average:\n%s\n", average)
  terminal.Refresh()

  return n.null
}

func (n *NodeRenderer) Name() string {
  return "renderer"
}

func (n *NodeRenderer) Rate() time.Duration {
  return n.hz
}

func (n *NodeRenderer) ValueLen() int64 {
  return 1
}

type NodeRender struct {
  p     *NodeRenderer
  node  Node
  polls int64
  val   string
  avg   string
}

func (r *NodeRender) Poll() bool {
  n := r.node
  old := r.polls
  r.polls = r.p.GetTracker().Value(n.Name()).Polls()
  if r.polls > old {
    r.getValue()
    return true
  }
  return false
}

func (r *NodeRender) getValue() {
  n := r.node
  u := n.Unit()
  tv := n.GetTracker().Value(n.Name())
  v := tv.Value()
  if v.ByteCapacity() == 0 {
    r.val = ""
    r.avg = ""
    return
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

  r.val = fmt.Sprintf(r.p.nameFmt + ": %s", n.Name(), vs)
  if as != "" {
    r.avg = fmt.Sprintf(r.p.nameFmt + ": %s (%.0fpps, %d)", n.Name(), as, math.Round(tv.PollsPerSecond()), tv.Polls())
  }
}
