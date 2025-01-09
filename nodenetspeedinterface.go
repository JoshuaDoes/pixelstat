package main

import (
  "github.com/JoshuaDoes/crunchio"

  "fmt"
  "sync"
  "time"
)

var (
  units = []string{"K", "M", "G", "T", "P", "Z"}
)

type NodeNetSpeedInterface struct {
  sync.Mutex
  *NodeFile

  iF, typ string
  bits bool

  b int64
  lastUnit string
  lastPoll time.Time
}

func NewNodeNetSpeedInterface(bits bool, iFPath, iF, typ string) *NodeNetSpeedInterface {
  n := new(NodeNetSpeedInterface)
  n.NodeFile = NewNodeFile(iFPath + "/" + iF + "/statistics/" + typ + "_bytes")
  n.iF = iF
  n.typ = typ
  n.bits = bits

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline
  n.b = strtoi64(v.String())
  n.lastPoll = time.Now()

  return n
}

func (n *NodeNetSpeedInterface) Name() string {
  return n.iF + " " + n.typ
}

func (n *NodeNetSpeedInterface) Rate() int {
  return 4
}

func (n *NodeNetSpeedInterface) Unit() string {
  return n.lastUnit
}

func (n *NodeNetSpeedInterface) ValueType() string {
  return "str"
}

func (n *NodeNetSpeedInterface) ValueLen() int64 {
  return 11 //999.99 PB/s or Pbps
}

func (n *NodeNetSpeedInterface) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  bytes := n.NodeFile.Value()
  bytes.TruncateRight(1) //Remove the newline
  b := strtoi64(bytes.String())

  str := "0 "
  if n.bits {
    str += "bps"
  } else {
    str += "B/s"
  }

  if b > n.b {
    dur := float64(time.Since(n.lastPoll).Nanoseconds()) / 1000 / 1000 / 1000
    n.lastPoll = time.Now()

    bytes := b - n.b
    n.b = b
    bps, unit := byteUnit(float64(bytes) / dur)

    if n.bits {
      bps *= 8
      unit += "bps"
    } else {
      unit += "B/s"
    }

    str = fmt.Sprintf("%.2f %s", bps, unit)
  }

  return crunchio.NewBuffer([]byte(str))
}

func byteUnit(bytes float64) (float64, string) {
  unit := ""
  for i := 0; i < len(units); i++ {
    if bytes < 1000 {
      break
    }
    bytes /= 1000
    unit = units[i]
  }
  return bytes, unit
}
