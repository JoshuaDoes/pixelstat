package main

import (
  "github.com/JoshuaDoes/crunchio"

  "fmt"
  "strings"
  "sync"
)

var (
  errCPUFreqNoCores = fmt.Errorf("cpufreqpolicy: no cores")
)

type NodeCPUFreqPolicy struct {
  sync.Mutex
  *NodeFile

  policy string
  cpus   []string
}

func NewNodeCPUFreqPolicy(policyPath, policy string) (*NodeCPUFreqPolicy, error) {
  n := new(NodeCPUFreqPolicy)
  n.policy = policy
  path := policyPath + "/" + policy
  ncpus, err := NewNodeFile(path + "/affected_cpus")
  if err != nil {
    return nil, errCPUFreqNoCores
  }
  n.cpus = strings.Split(string(ncpus.Value().Bytes()), " ")
  nf, err := NewNodeFile(path + "/cpuinfo_cur_freq")
  if err != nil {
    nf, err = NewNodeFile(path + "/scaling_cur_freq")
  }
  n.NodeFile = nf
  return n, err
}

func (n *NodeCPUFreqPolicy) Name() string {
  return n.policy + " freq"
}

func (n *NodeCPUFreqPolicy) Rate() int {
  return 1000
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
  if v == nil || v.ByteCapacity() == 0 {
    return nil
  }
  v.TruncateRight(1) //Remove the newline

  Hz := strtoi32(v.String())
  MHz := Hz / 1000

  v = crunchio.NewBuffer()
  v.Grow(4)
  v.WriteAbstract(MHz)
  v.Seek(0, 0)
  return v
}

func (n *NodeCPUFreqPolicy) CPUs() []string {
  return n.cpus
}
