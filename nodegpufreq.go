package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
)

type NodeGPUFreq struct {
  sync.Mutex
  *NodeFile
}

func NewNodeGPUFreq(path string) (*NodeGPUFreq, error) {
  n := new(NodeGPUFreq)
  nf, err := NewNodeFile(path)
  n.NodeFile = nf
  return n, err
}

func (n *NodeGPUFreq) Name() string {
  return "gpu"
}

func (n *NodeGPUFreq) Rate() int {
  return 1000
}

func (n *NodeGPUFreq) Unit() string {
  return " MHz"
}

func (n *NodeGPUFreq) ValueType() string {
  return "i32"
}

func (n *NodeGPUFreq) ValueLen() int64 {
  return 1
}

func (n *NodeGPUFreq) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  if v == nil || v.ByteCapacity() == 0 {
    return nil
  }
  v.TruncateRight(1) //Remove the newline

  Hz := strtoi32(v.String())
  MHz := Hz / 1000 / 1000

  v = crunchio.NewBuffer()
  v.Grow(4)
  v.WriteAbstract(MHz)
  v.Seek(0, 0)
  return v
}
