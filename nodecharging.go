package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
)

type NodeCharging struct {
  sync.Mutex
  *NodeFile
}

func NewNodeCharging(path string) (*NodeCharging, error) {
  n := new(NodeCharging)
  nf, err := NewNodeFile(path)
  n.NodeFile = nf
  return n, err
}

func (n *NodeCharging) Name() string {
  return "charging"
}

func (n *NodeCharging) Rate() int {
  return 2
}

func (n *NodeCharging) Unit() string {
  return ""
}

func (n *NodeCharging) ValueType() string {
  return "str"
}

func (n *NodeCharging) ValueLen() int64 {
  return 12
}
func (n *NodeCharging) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()
  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline
  v.Seek(0, 0)
  return v
}
