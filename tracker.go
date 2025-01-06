package main

import (
  "github.com/JoshuaDoes/crunchio"

  "fmt"
  "sync"
)

var (
  errNodeNameless = fmt.Errorf("tracker: cannot register nameless node")
  errNodeExists   = fmt.Errorf("tracker: cannot duplicate root node")
)

type Tracker struct {
  nodes []Node
  vals  map[string]*TrackerValue
}

func NewTracker(nodes ...Node) *Tracker {
  t := new(Tracker)
  t.nodes = make([]Node, 0)
  t.vals = make(map[string]*TrackerValue)
  for i := 0; i < len(nodes); i++ {
    t.Register(nodes[i])
  }
  return t
}

func (t *Tracker) Start() {
  for _, tv := range t.vals {
    tv.start()
  }
}

func (t *Tracker) Close() {
  for _, tv := range t.vals {
    tv.close()
  }
  t.vals = nil
  t.nodes = nil
}

func (t *Tracker) Nodes() []Node {
  nodes := make([]Node, 0)
  for i := 0; i < len(t.nodes); i++ {
    name := t.nodes[i].Name()
    if _, exists := t.vals[name]; exists {
      nodes = append(nodes, t.nodes[i])
    }
  }
  return nodes
}

func (t *Tracker) Node(name string) Node {
  for i := 0; i < len(t.nodes); i++ {
    if t.nodes[i].Name() == name {
      return t.nodes[i]
    }
  }
  return nil
}

func (t *Tracker) Value(node string) *crunchio.Buffer {
  return t.vals[node].Value()
}

func (t *Tracker) Register(n Node) {
  name := n.Name()
  if name == "" {
    perr(errNodeNameless)
  }
  if t.Node(name) != nil {
    perr(errNodeExists)
  }
  n.SetTracker(t)
  t.nodes = append(t.nodes, n)

  if n.ValueLen() > 0 {
    tv := new(TrackerValue)
    tv.node = n
    t.vals[name] = tv
  }
}

type TrackerValue struct {
  sync.Mutex
  node    Node
  value   *crunchio.Buffer
  stopper chan bool
}

func (tv *TrackerValue) start() {
  if tv.stopper != nil {
    return
  }

  if tv.node.Rate() > 0 {
    tv.stopper = loop(tv.node.Rate(), tv.getValue)
  }
}

func (tv *TrackerValue) close() {
  tv.stopper <- true
  tv.stopper = nil
}

func (tv *TrackerValue) getValue() {
  tv.Lock()
  defer tv.Unlock()
  tv.value = tv.node.Value()
}

func (tv *TrackerValue) Value() *crunchio.Buffer {
  if tv.value == nil {
    tv.getValue()
  }
  return tv.value
}
