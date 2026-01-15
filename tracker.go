package main

import (
	"fmt"
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

func (t *Tracker) Start(maxHz int) {
	vals := make([]*TrackerValue, 0)
	for i := 0; i < len(t.nodes); i++ {
		t.nodes[i].SetTracker(t)
		name := t.nodes[i].Name()
		if tv, exists := t.vals[name]; exists {
			vals = append(vals, tv)
		}
	}
	for i := 0; i < len(vals); i++ {
		go vals[i].start(maxHz)
	}
}

func (t *Tracker) Close() {
	for _, tv := range t.vals {
		tv.close()
	}
	t.vals = nil
	t.nodes = nil
}

func (t *Tracker) Closed() bool {
	return t.vals == nil && t.nodes == nil
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

func (t *Tracker) Value(node string) *TrackerValue {
	return t.vals[node]
}

func (t *Tracker) Register(nodes ...Node) {
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		name := n.Name()
		if name == "" {
			perr(errNodeNameless)
		}
		if t.Node(name) != nil {
			perr(errNodeExists)
		}
		t.nodes = append(t.nodes, n)

		if n.ValueLen() > 0 {
			tv := new(TrackerValue)
			tv.node = n
			t.vals[name] = tv
		}
	}
}
