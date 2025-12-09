package main

import (
  "github.com/JoshuaDoes/crunchio"

  "fmt"
  "sync"
)

var (
  errCPUBandwidthNeedFreqPol = fmt.Errorf("wattage: missing node cpufreqpolicy")
)

type NodeCPUBandwidth struct {
  sync.Mutex
  *Tracker

  rate int
  pols []Node
}

func (n *NodeCPUBandwidth) SetTracker(t *Tracker) {
  n.pols = make([]Node, 0)
  for _, pol := range t.Nodes() {
    name := pol.Name()
    if len(name) < 7 {
      continue
    }
    if name[:6] == "policy" {
      n.pols = append(n.pols, pol)
      rate := pol.Rate()
      if rate > n.rate {
        n.rate = rate
      }
    }
  }
  n.Tracker = t
}

func (n *NodeCPUBandwidth) GetTracker() *Tracker {
  return n.Tracker
}

func NewNodeCPUBandwidth() *NodeCPUBandwidth {
  n := new(NodeCPUBandwidth)
  return n
}

func (n *NodeCPUBandwidth) Name() string {
  return "cpu"
}

func (n *NodeCPUBandwidth) Rate() int {
  return n.rate
}

func (n *NodeCPUBandwidth) Unit() string {
  return " MHz"
}

func (n *NodeCPUBandwidth) ValueType() string {
  return "i32"
}

func (n *NodeCPUBandwidth) ValueLen() int64 {
  return 1
}

func (n *NodeCPUBandwidth) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  bandwidth := int32(0)
  v := crunchio.NewBuffer()
  v.Grow(4)
  for _, pol := range n.pols {
    if pol == nil {
      continue //TODO: Drop from n.pols
    }
    value := pol.Value()
    if value == nil {
      continue //TODO: Drop from n.pols
    }
    bw := value.ReadI32LENext(1)
    if len(bw) == 0 {
      continue
    }
    npol := pol.(*NodeCPUFreqPolicy)
    bandwidth += int32(len(npol.CPUs())) * bw[0]
  }
  v.WriteAbstract(bandwidth)
  v.Seek(0, 0)
  return v
}

func (n *NodeCPUBandwidth) Close() error {
  n.Lock()
  defer n.Unlock()
  n.pols = nil
  n.rate = 0
  return nil
}
