package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
  "time"
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

func (n *NodeCurrent) Rate() time.Duration {
  return hertz(4)
}

func (n *NodeCurrent) Unit() string {
  return " mA"
}

func (n *NodeCurrent) ValueLen() int {
  return 5
}

func (n *NodeCurrent) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline

  microA := str2float(v.String())
  milliA := microA / 1000

  v = crunchio.NewBuffer([]byte(float2str(milliA, 2)))
  return v
}
