package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
  "time"
)

type NodeCPUFreqPolicy struct {
  sync.Mutex
  *NodeFile

  policy string
}

func NewNodeCPUFreqPolicy(policyPath, policy string) *NodeCPUFreqPolicy {
  n := new(NodeCPUFreqPolicy)
  n.policy = policy
  n.NodeFile = NewNodeFile(policyPath + "/" + policy + "/cpuinfo_cur_freq")
  return n
}

func (n *NodeCPUFreqPolicy) Name() string {
  return n.policy + " freq"
}

func (n *NodeCPUFreqPolicy) Rate() time.Duration {
  return hertz(240)
}

func (n *NodeCPUFreqPolicy) Unit() string {
  return " MHz"
}

func (n *NodeCPUFreqPolicy) ValueType() string {
  return "i32"
}

func (n *NodeCPUFreqPolicy) ValueLen() int64 {
  return 1
}

func (n *NodeCPUFreqPolicy) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline

  Hz := strtoi32(v.String())
  MHz := Hz / 1000

  v = crunchio.NewBuffer()
  v.Grow(4)
  v.WriteAbstract(MHz)
  v.Seek(0, 0)
  return v
}
