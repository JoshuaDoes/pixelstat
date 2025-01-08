package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
  "time"
)

type NodeBattery struct {
  sync.Mutex
  *NodeFile
}

func NewNodeBattery(path string) *NodeBattery {
  n := new(NodeBattery)
  n.NodeFile = NewNodeFile(path)
  return n
}

func (n *NodeBattery) Name() string {
  return "battery"
}

func (n *NodeBattery) Rate() time.Duration {
  return hertz(2)
}

func (n *NodeBattery) Unit() string {
  return "%"
}

func (n *NodeBattery) ValueType() string {
  return "str"
}

func (n *NodeBattery) ValueLen() int64 {
  return 3
}

func (n *NodeBattery) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()
  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline
  v.Seek(0, 0)
  return v
}
