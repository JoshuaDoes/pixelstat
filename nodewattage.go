package main

import (
  "github.com/JoshuaDoes/crunchio"

  "fmt"
  "sync"
  "time"
)

var (
  errWattageNeedCurrent = fmt.Errorf("wattage: missing node current")
  errWattageNeedVoltage = fmt.Errorf("wattage: missing node voltage")
)

type NodeWattage struct {
  sync.Mutex
  *Tracker
}

func (n *NodeWattage) SetTracker(t *Tracker) {
  if t.Node("current") == nil {
    perr(errWattageNeedCurrent)
  }
  if t.Node("voltage") == nil {
    perr(errWattageNeedVoltage)
  }
  n.Tracker = t
}

func (n *NodeWattage) GetTracker() *Tracker {
  return n.Tracker
}

func NewNodeWattage() *NodeWattage {
  n := new(NodeWattage)
  return n
}

func (n *NodeWattage) Name() string {
  return "wattage"
}

func (n *NodeWattage) Rate() time.Duration {
  return hertz(5)
}

func (n *NodeWattage) Unit() string {
  return " W"
}

func (n *NodeWattage) ValueType() string {
  return "f64"
}

func (n *NodeWattage) ValueLen() int64 {
  return 1
}

func (n *NodeWattage) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  amps := n.GetTracker().Value("current").Value().ReadF64LENext(1)[0] / 1000
  volts := n.GetTracker().Value("voltage").Value().ReadF64LENext(1)[0] / 1000

  wattage := amps * volts

  v := crunchio.NewBuffer()
  v.Grow(8)
  v.WriteAbstract(wattage)
  v.Seek(0, 0)
  return v
}
