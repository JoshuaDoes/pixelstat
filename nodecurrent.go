package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
)

type NodeCurrent struct {
  sync.Mutex
  *NodeFile
}

func NewNodeCurrent(path string) *NodeCurrent {
  n := new(NodeCurrent)
  n.NodeFile = NewNodeFile(path)
  return n
}

func (n *NodeCurrent) Name() string {
  return "current"
}

func (n *NodeCurrent) Rate() int {
  return 5
}

func (n *NodeCurrent) Unit() string {
  return " mA"
}

func (n *NodeCurrent) ValueType() string {
  return "f64"
}

func (n *NodeCurrent) ValueLen() int64 {
  return 1
}

func (n *NodeCurrent) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline

  microA := strtof64(v.String())
  milliA := microA / 1000

  v = crunchio.NewBuffer()
  v.Grow(8)
  v.WriteAbstract(milliA)
  v.Seek(0, 0)
  return v
}
