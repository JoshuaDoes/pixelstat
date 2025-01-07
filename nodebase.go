package main

import (
  "github.com/JoshuaDoes/crunchio"

  "time"
)

type NodeBase struct {
  *Tracker
}

func NewNodeBase() *NodeBase {
  return new(NodeBase)
}

func (n *NodeBase) SetTracker(t *Tracker) {
  n.Tracker = t
}

func (n *NodeBase) GetTracker() *Tracker {
  return n.Tracker
}

func (n *NodeBase) Name() string {
  return "base"
}

func (n *NodeBase) Rate() time.Duration {
  return 0
}

func (n *NodeBase) Unit() string {
  return ""
}

func (n *NodeBase) ValueType() string {
  return ""
}

func (n *NodeBase) ValueLen() int64 {
  return 0
}

func (n *NodeBase) Value() *crunchio.Buffer {
  return crunchio.NewBuffer(make([]byte, 0))
}
