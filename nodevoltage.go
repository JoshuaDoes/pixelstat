package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
)

type NodeVoltage struct {
  sync.Mutex
  *NodeFile
}

func NewNodeVoltage(path string) (*NodeVoltage, error) {
  n := new(NodeVoltage)
  nf, err := NewNodeFile(path)
  n.NodeFile = nf
  return n, err
}

func (n *NodeVoltage) Name() string {
  return "voltage"
}

func (n *NodeVoltage) Rate() int {
  return 5
}

func (n *NodeVoltage) Unit() string {
  return " mV"
}

func (n *NodeVoltage) ValueType() string {
  return "f64"
}

func (n *NodeVoltage) ValueLen() int64 {
  return 1
}

func (n *NodeVoltage) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  if v == nil || v.Size() < 2 {
    return nil
  }
  v.TruncateRight(1) //Remove the newline

  microV := strtof64(v.String())
  milliV := microV / 1000

  v = crunchio.NewBuffer()
  v.Grow(8)
  v.WriteAbstract(milliV)
  v.Seek(0, 0)
  return v
}
