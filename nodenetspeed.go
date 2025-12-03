package main

import (
  "fmt"
  "io/fs"
  "os"
)

var (
  errNetNotSymlink = fmt.Errorf("net: net path contains non-symlinks")
)

type NodeNetSpeed struct {
  *NodeBase
  iFPath string
  iFs []string
  bits bool
}

func (n *NodeNetSpeed) SetTracker(t *Tracker) {
  dir, err := os.ReadDir(n.iFPath)
  if err != nil {
    return
  }

  iFs := make([]string, 0)
  for i := 0; i < len(dir); i++ {
    file := dir[i]
    if file.Type() != fs.ModeSymlink {
      continue
    }

    iF := file.Name()
    found := true
    if len(n.iFs) > 0 {
      found = false
      for j := 0; j < len(n.iFs); j++ {
        if n.iFs[j] == iF {
          found = true
          break
        }
      }
    }
    if !found {
      continue
    }

    rxFile, err := os.Open(n.iFPath + "/" + iF + "/statistics/rx_bytes")
    if err != nil {
      continue
    }
    rxFile.Close()
    txFile, err := os.Open(n.iFPath + "/" + iF + "/statistics/tx_bytes")
    if err != nil {
      continue
    }
    txFile.Close()

    iFs = append(iFs, iF)
  }

  for i := 0; i < len(iFs); i++ {
    node, err := NewNodeNetSpeedInterface(n.bits, n.iFPath, iFs[i], "rx")
    if err == nil {
      t.Register(node)
    }
    node, err = NewNodeNetSpeedInterface(n.bits, n.iFPath, iFs[i], "tx")
    if err == nil {
      t.Register(node)
    }
  }
}

func NewNodeNetSpeed(bits bool, iFPath string, iFs ...string) *NodeNetSpeed {
  n := new(NodeNetSpeed)
  n.NodeBase = NewNodeBase()
  n.iFPath = iFPath
  n.iFs = iFs
  n.bits = bits
  return n
}

func (n *NodeNetSpeed) Name() string {
  return "netspeed"
}
