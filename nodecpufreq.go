package main

import (
  "os"
)

type NodeCPUFreq struct {
  *NodeBase
  policyPath string
}

func (n *NodeCPUFreq) SetTracker(t *Tracker) {
  dir, err := os.ReadDir(n.policyPath)
  if err != nil {
    return
  }

  policies := make([]string, 0)
  for _, file := range dir {
    if !file.IsDir() {
      continue
    }
    policy := file.Name()
    freqFile, err := os.Open(n.policyPath + "/" + policy + "/cpuinfo_cur_freq")
    if err != nil {
      continue
    }
    freqFile.Close()
    policies = append(policies, policy)
  }

  if len(policies) == 0 {
    return
  }

  for i := 0; i < len(policies); i++ {
    node, err := NewNodeCPUFreqPolicy(n.policyPath, policies[i])
    if err == nil {
      t.Register(node)
    }
  }
}

func NewNodeCPUFreq(policyPath string) *NodeCPUFreq {
  n := new(NodeCPUFreq)
  n.NodeBase = NewNodeBase()
  n.policyPath = policyPath
  return n
}

func (n *NodeCPUFreq) Name() string {
  return "cpufreq"
}
