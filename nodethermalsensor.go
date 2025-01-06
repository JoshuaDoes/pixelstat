package main

import (
  "github.com/JoshuaDoes/crunchio"

  "os"
  "sync"
  "time"
)

type NodeThermalSensor struct {
  sync.Mutex
  *NodeFile

  sensor string
  ttype  string
}

func NewNodeThermalSensor(thermalPath, sensor string) *NodeThermalSensor {
  data, err := os.ReadFile(thermalPath + "/" + sensor + "/type")
  perr(err)
  ttype := string(data[:len(data)-1])
  
  n := new(NodeThermalSensor)
  n.sensor = sensor
  n.ttype = ttype
  n.NodeFile = NewNodeFile(thermalPath + "/" + sensor + "/temp")
  return n
}

func (n *NodeThermalSensor) Name() string {
  return n.ttype + " temp"
}

func (n *NodeThermalSensor) Rate() time.Duration {
  return hertz(5)
}

func (n *NodeThermalSensor) Unit() string {
  return "°C"
}

func (n *NodeThermalSensor) ValueLen() int {
  return 6
}

func (n *NodeThermalSensor) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()

  v := n.NodeFile.Value()
  v.TruncateRight(1) //Remove the newline

  c := str2float(v.String()) / 1000

  v = crunchio.NewBuffer([]byte(float2str(c, 2)))
  return v
}
