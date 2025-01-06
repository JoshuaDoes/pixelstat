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
  return hertz(4)
}

func (n *NodeWattage) Unit() string {
  return " W"
}

func (n *NodeWattage) ValueLen() int {
  return 7
}

func (n *NodeWattage) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  amps := str2float(n.GetTracker().Value("current").String()) / 1000
  volts := str2float(n.GetTracker().Value("voltage").String()) / 1000

  wattage := amps * volts

  v := crunchio.NewBuffer([]byte(float2str(wattage, 2)))
  return v
}
