package main

import (
  "github.com/JoshuaDoes/crunchio"

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
  null string

  b int64
  lastPoll time.Time
}

func NewNodeNetSpeedInterface(bits bool, iFPath, iF, typ string) (*NodeNetSpeedInterface, error) {
  n := new(NodeNetSpeedInterface)
  nf, err := NewNodeFile(iFPath + "/" + iF + "/statistics/" + typ + "_bytes")
  if err != nil {
    return nil, err
  }
  n.NodeFile = nf
  n.iF = iF
  n.typ = typ
  n.bits = bits

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline
  n.b = strtoi64(v.String())
  n.lastPoll = time.Now()

  n.null = "0 "
  if bits {
    n.null += "bps"
  } else {
    n.null += "B/s"
  }

  return n, nil
}

func (n *NodeNetSpeedInterface) Name() string {
  return n.iF + " " + n.typ
}

func (n *NodeNetSpeedInterface) Rate() int {
  return 2
}

func (n *NodeNetSpeedInterface) Unit() string {
  return ""
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

  str := n.null

  v := n.NodeFile.Value()
  if v == nil || v.ByteCapacity() == 0 {
    return nil
  }
  v.TruncateRight(1) //Remove the newline
  b := strtoi64(v.String())

  if b > n.b {
    dur := float64(time.Since(n.lastPoll).Nanoseconds()) / 1000 / 1000 / 1000
    n.lastPoll = time.Now()

    now := b - n.b
    if n.bits {
      now *= 8
    }
    n.b = b

    bps, unit := byteUnit(float64(now) / dur)
    if n.bits {
      unit += "bps"
    } else {
      unit += "B/s"
    }

    str = f64tostr(bps, 2) + " " + unit
  }

  pad := int(n.ValueLen()) - len(str)
  for i := 0; i < pad; i++ {
    str = " " + str
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
