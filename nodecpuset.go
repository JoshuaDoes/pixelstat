package main

import (
  "os"
)

type NodeCPUSet struct {
  *NodeBase
  setsPath string
}

func (n *NodeCPUSet) SetTracker(t *Tracker) {
  dir, err := os.ReadDir(n.setsPath)
  if err != nil {
    return
  }

  sets := make([]string, 0)
  for i := 0; i < len(dir); i++ {
    file := dir[i]
    if !file.IsDir() {
      continue
    }

    set := file.Name()

    cpusFile, err := os.Open(n.setsPath + "/" + set + "/cpus")
    perr(err)
    cpusFile.Close()
    tasksFile, err := os.Open(n.setsPath + "/" + set + "/tasks")
    perr(err)
    tasksFile.Close()

    sets = append(sets, set)
  }

  for i := 0; i < len(sets); i++ {
    node, err := NewNodeCPUSetTasks(n.setsPath, sets[i])
    if err == nil {
      t.Register(node)
    }
  }
}

func NewNodeCPUSet(setsPath string) *NodeCPUSet {
  n := new(NodeCPUSet)
  n.NodeBase = NewNodeBase()
  n.setsPath = setsPath
  return n
}

func (n *NodeCPUSet) Name() string {
  return "cpuset"
}
