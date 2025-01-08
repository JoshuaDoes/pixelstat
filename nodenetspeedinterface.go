package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
  "time"
)

type NodeNetSpeedInterface struct {
  sync.Mutex
  *NodeFile

  iF, typ string
  b int64
}

func NewNodeNetSpeedInterface(iFPath, iF, typ string) *NodeNetSpeedInterface {
  n := new(NodeNetSpeedInterface)
  n.iF = iF
  n.typ = typ
  n.NodeFile = NewNodeFile(iFPath + "/" + iF + "/statistics/" + typ + "_bytes")
  return n
}

func (n *NodeNetSpeedInterface) Name() string {
  return n.iF + " " + n.typ
}

func (n *NodeNetSpeedInterface) Rate() time.Duration {
  return hertz(1)
}

func (n *NodeNetSpeedInterface) Unit() string {
  return " KB/s"
}

func (n *NodeNetSpeedInterface) ValueType() string {
  return "f64"
}

func (n *NodeNetSpeedInterface) ValueLen() int64 {
  return 1
}

func (n *NodeNetSpeedInterface) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline

  b := strtoi64(v.String())
  kbps := float64(0)
  if b > n.b {
    last := b - n.b
    kbps = (float64(n.Rate()) / float64(last)) / 1000
    n.b = b
  }

  v = crunchio.NewBuffer()
  v.Grow(8)
  v.WriteAbstract(kbps)
  v.Seek(0, 0)
  return v
}
