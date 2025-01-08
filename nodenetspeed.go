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
}

func (n *NodeNetSpeed) SetTracker(t *Tracker) {
  dir, err := os.ReadDir(n.iFPath)
  perr(err)

  iFs := make([]string, 0)
  for i := 0; i < len(dir); i++ {
    file := dir[i]
    if file.Type() != fs.ModeSymlink {
      perr(errNetNotSymlink)
    }

    iF := file.Name()
    found := false
    for j := 0; j < len(n.iFs); j++ {
      if n.iFs[j] == iF {
        found = true
        break
      }
    }
    if !found {
      continue
    }

    rxFile, err := os.Open(n.iFPath + "/" + iF + "/statistics/rx_bytes")
    perr(err)
    rxFile.Close()
    txFile, err := os.Open(n.iFPath + "/" + iF + "/statistics/tx_bytes")
    perr(err)
    txFile.Close()

    iFs = append(iFs, iF)
  }

  for i := 0; i < len(iFs); i++ {
    t.Register(NewNodeNetSpeedInterface(n.iFPath, iFs[i], "rx"))
    t.Register(NewNodeNetSpeedInterface(n.iFPath, iFs[i], "tx"))
  }
}

func NewNodeNetSpeed(iFPath string, iFs ...string) *NodeNetSpeed {
  n := new(NodeNetSpeed)
  n.NodeBase = NewNodeBase()
  n.iFPath = iFPath
  n.iFs = iFs
  return n
}

func (n *NodeNetSpeed) Name() string {
  return "netspeed"
}
