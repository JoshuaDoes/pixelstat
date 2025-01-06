package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
  "time"
)

type NodeCharging struct {
  sync.Mutex
  *NodeFile
}

func NewNodeCharging(path string) *NodeCharging {
  n := new(NodeCharging)
  n.NodeFile = NewNodeFile(path)
  return n
}

func (n *NodeCharging) Name() string {
  return "charging"
}

func (n *NodeCharging) Rate() time.Duration {
  return hertz(2)
}

func (n *NodeCharging) Unit() string {
  return ""
}

func (n *NodeCharging) ValueLen() int {
  return 11
}
func (n *NodeCharging) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()
  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline
  return v
}
