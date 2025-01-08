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

  polls      []*NodeRender
}

func NewNodeRenderer(hz int, vrr bool) *NodeRenderer {
  n := new(NodeRenderer)
  n.NodeBase = NewNodeBase()
  n.hz = hertz(time.Duration(hz))
  n.vrr = vrr
  n.null = crunchio.NewBuffer(make([]byte, 1))
  n.polls = make([]*NodeRender, 0)
  terminal = ncurses.Init()
  return n
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

    r := new(NodeRender)
    r.p = n
    r.node = node
    r.tv = t.Value(node.Name())
    n.polls = append(n.polls, r)
  }

  n.nameLen = nameLen
  n.nameFmt = "%" + fmt.Sprintf("%d", nameLen) + "s"
  n.unitLen = unitLen
  n.Tracker = t
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

func (n *NodeRenderer) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()
  if terminal == nil {
    return nil
  }
  pollers := n.Pollers()

  go func(n *NodeRenderer) {
    for i := 0; i < pollers; i++ {
      go n.GetRender(i).Poll()
    }
  }(n)

  if n.vrr {
    newFrame := false
    for i := 0; i < pollers; i++ {
      if n.GetRender(i).NewFrame() {
        newFrame = true
        break
      }
    }
    if !newFrame {
      return nil
    }
  }

  now := ""
  average := ""
  for i := 0; i < pollers; i++ {
    p := n.GetRender(i)
    if val := p.val; val != "" {
      now += val + "\n"
    }
    if avg := p.avg; avg != "" {
      average += avg + "\n"
    }
  }

  if terminal == nil {
    return nil
  }
  rendertv := n.GetTracker().Value(n.Name())
  terminal.Erase()
  terminal.Printf("%.0f FPS\n\n", math.Round(rendertv.PollsPerSecond()))
  terminal.Printf("Now:\n%s\n", now)
  terminal.Printf("Average:\n%s\n", average)
  terminal.Refresh()

  return n.null
}

func (n *NodeRenderer) Close() error {
  ncurses.EndWin()
  terminal = nil
  return nil
}

func (n *NodeRenderer) Pollers() int {
  return len(n.polls)
}

func (n *NodeRenderer) GetRender(i int) *NodeRender {
  return n.polls[i]
}

type NodeRender struct {
  sync.Mutex

  p     *NodeRenderer
  node  Node
  tv    *TrackerValue
  polls int64
  val   string
  avg   string
  newF  bool
  poll  bool
}

func (r *NodeRender) Poll() {
  if r.poll {
    return
  }
  r.Lock()
  defer r.Unlock()
  r.poll = true

  n := r.node
  old := r.polls
  r.polls = r.p.GetTracker().Value(n.Name()).Polls()
  if r.polls > old {
    r.getValue()
    r.newF = true
  }
  r.poll = false
}

func (r *NodeRender) NewFrame() bool {
  r.Lock()
  defer r.Unlock()

  newF := r.newF
  r.newF = false
  return newF
}

func (r *NodeRender) getValue() {
  v := r.tv.Value()
  if v == nil || v.ByteCapacity() == 0 {
    r.val = ""
    r.avg = ""
    return
  }

  vs := ""
  as := ""
  vt := r.node.ValueType()
  switch vt {
  case "raw":
    length := r.node.ValueLen()
    if bc := v.ByteCapacity(); bc < length {
      length = bc
    }
    vs = fmt.Sprintf("% X", v.ReadBytes(0, length)) + r.node.Unit()
  case "str":
    length := r.node.ValueLen()
    if bc := v.ByteCapacity(); bc < length {
      length = bc
    }
    vs = string(v.ReadBytes(0, length)) + r.node.Unit()
  case "i32":
    length := r.node.ValueLen()
    if bc := v.ByteCapacity(); (bc / 4) < length {
      length = (bc / 4)
    }
    for i := int64(0); i < length; i++ {
      if i > 0 {
        vs += " "
        as += " "
      }
      vs += i32tostr(v.ReadI32LENext(1)[0]) + r.node.Unit()
      as += f64tostr(r.tv.Average(i), 0) + r.node.Unit()
    }
  case "i64":
    length := r.node.ValueLen()
    if bc := v.ByteCapacity(); (bc / 8) < length {
      length = (bc / 8)
    }
    for i := int64(0); i < length; i++ {
      if i > 0 {
        vs += " "
        as += " "
      }
      vs += i64tostr(v.ReadI64LENext(1)[0]) + r.node.Unit()
      as += f64tostr(r.tv.Average(i), 0) + r.node.Unit()
    }
  case "f32":
    length := r.node.ValueLen()
    if bc := v.ByteCapacity(); (bc / 4) < length {
      length = (bc / 4)
    }
    for i := int64(0); i < length; i++ {
      if i > 0 {
        vs += " "
        as += " "
      }
      vs += f32tostr(v.ReadF32LENext(1)[0], 2) + r.node.Unit()
      as += f64tostr(r.tv.Average(i), 2) + r.node.Unit()
    }
  case "f64":
    length := r.node.ValueLen()
    if bc := v.ByteCapacity(); (bc / 8) < length {
      length = (bc / 8)
    }
    for i := int64(0); i < length; i++ {
      if i > 0 {
        vs += " "
        as += " "
      }
      vs += f64tostr(v.ReadF64LENext(1)[0], 2) + r.node.Unit()
      as += f64tostr(r.tv.Average(i), 2) + r.node.Unit()
    }
  default:
    perr(fmt.Errorf("renderer: invalid type: %s", vt))
  }

  r.val = fmt.Sprintf(r.p.nameFmt + ": %s", r.node.Name(), vs)
  if as != "" {
    r.avg = fmt.Sprintf(r.p.nameFmt + ": %s (%.0fpps, %d)", r.node.Name(), as, math.Round(r.tv.PollsPerSecond()), r.tv.Polls())
  }
}
