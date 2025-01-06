package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
  "time"
)

type NodeVoltage struct {
  sync.Mutex
  *NodeFile
}

func NewNodeVoltage(path string) *NodeVoltage {
  n := new(NodeVoltage)
  n.NodeFile = NewNodeFile(path)
  return n
}

func (n *NodeVoltage) Name() string {
  return "voltage"
}

func (n *NodeVoltage) Rate() time.Duration {
  return hertz(4)
}

func (n *NodeVoltage) Unit() string {
  return " mV"
}

func (n *NodeVoltage) ValueLen() int {
  return 5
}

func (n *NodeVoltage) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline

  microV := str2int64(v.String())
  milliV := microV / 1000

  v = crunchio.NewBuffer([]byte(int642str(milliV)))
  return v
}
