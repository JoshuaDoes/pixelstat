package main

import (
  "fmt"
  "os"
)

var (
  errCPUFreqPolicies     = fmt.Errorf("cpufreq: no policies were found")
  errCPUFreqPolicyNotDir = fmt.Errorf("cpufreq: policy path contains files")
)

type NodeCPUFreq struct {
  *NodeBase
  policyPath string
}

func (n *NodeCPUFreq) SetTracker(t *Tracker) {
  dir, err := os.ReadDir(n.policyPath)
  perr(err)

  policies := make([]string, 0)
  for _, file := range dir {
    if !file.IsDir() {
      perr(errCPUFreqPolicyNotDir)
    }
    policy := file.Name()
    freqFile, err := os.Open(n.policyPath + "/" + policy + "/cpuinfo_cur_freq")
    perr(err)
    freqFile.Close()
    policies = append(policies, policy)
  }

  if len(policies) == 0 {
    perr(errCPUFreqPolicies)
  }

  for i := 0; i < len(policies); i++ {
    t.Register(NewNodeCPUFreqPolicy(n.policyPath, policies[i]))
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
