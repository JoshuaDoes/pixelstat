package main

import (
  "github.com/JoshuaDoes/crunchio"

  "fmt"
  "strings"
  "sync"
)

type NodeCPUSetTasks struct {
  sync.Mutex
  *NodeFile
  cpus *NodeFile

  set string
  cores string
}

func NewNodeCPUSetTasks(setsPath, set string) *NodeCPUSetTasks {
  n := new(NodeCPUSetTasks)
  n.NodeFile = NewNodeFile(setsPath + "/" + set + "/tasks")
  n.cpus = NewNodeFile(setsPath + "/" + set + "/cpus")
  n.set = set

  n.setCpus()

  return n
}

func (n *NodeCPUSetTasks) Name() string {
  return fmt.Sprintf("%s (%s)", n.set, n.cores)
}

func (n *NodeCPUSetTasks) Rate() int {
  return 240
}

func (n *NodeCPUSetTasks) Unit() string {
  return " tasks"
}

func (n *NodeCPUSetTasks) ValueType() string {
  return "i32"
}

func (n *NodeCPUSetTasks) ValueLen() int64 {
  return 1
}

func (n *NodeCPUSetTasks) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  n.setCpus()

  v := n.NodeFile.Value()
  if v == nil || v.ByteCapacity() == 0 {
    v = crunchio.NewBuffer()
    v.Grow(4)
    return v
  }
  v.TruncateRight(1) //Remove the newline

  tasks := strings.Split(v.String(), "\n")
  total := len(tasks)

  v = crunchio.NewBuffer()
  v.Grow(4)
  v.WriteAbstract(int32(total))
  v.Seek(0, 0)
  return v
}

func (n *NodeCPUSetTasks) setCpus() {
  cpus := n.cpus.Value()
  if cpus.ByteCapacity() > 0 {
    cpus.TruncateRight(1) //Remove the newline
  } else {
    cpus.Grow(1)
    cpus.WriteAbstract("?")
  }
  n.cores = cpus.String()
}
