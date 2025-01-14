package main

import (
  "github.com/JoshuaDoes/crunchio"

  "os"
  "sync"
)

type NodeThermalSensor struct {
  sync.Mutex
  *NodeFile

  sensor string
  ttype  string
}

func NewNodeThermalSensor(thermalPath, sensor string) (*NodeThermalSensor, error) {
  data, err := os.ReadFile(thermalPath + "/" + sensor + "/type")
  if err != nil {
    return nil, err
  }
  ttype := string(data[:len(data)-1])

  n := new(NodeThermalSensor)
  n.sensor = sensor
  n.ttype = ttype
  nf, err := NewNodeFile(thermalPath + "/" + sensor + "/temp")
  n.NodeFile = nf
  return n, err
}

func (n *NodeThermalSensor) Name() string {
  return n.ttype + " temp"
}

func (n *NodeThermalSensor) Rate() int {
  return 8
}

func (n *NodeThermalSensor) Unit() string {
  return "°C"
}

func (n *NodeThermalSensor) ValueType() string {
  return "f64"
}

func (n *NodeThermalSensor) ValueLen() int64 {
  return 1
}

func (n *NodeThermalSensor) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  if v == nil || v.ByteCapacity() == 0 {
    return nil
  }
  v.TruncateRight(1) //Remove the newline

  temp := strtof64(v.String()) / 1000

  v = crunchio.NewBuffer()
  v.Grow(8)
  v.WriteAbstract(temp)
  v.Seek(0, 0)
  return v
}
