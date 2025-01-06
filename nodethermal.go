package main

import (
  "fmt"
  "io/fs"
  "os"
)

var (
  errThermals           = fmt.Errorf("thermal: no thermals were found")
  errThermalsNotSymlink = fmt.Errorf("thermal: thermals path contains non-symlinks")
)

type NodeThermal struct {
  *NodeBase
  thermalPath string
}

func (n *NodeThermal) SetTracker(t *Tracker) {
  dir, err := os.ReadDir(n.thermalPath)
  perr(err)

  sensors := make([]string, 0)
  for _, file := range dir {
    if file.Type() != fs.ModeSymlink {
      perr(errThermalsNotSymlink)
    }
    sensor := file.Name()

    _, err := os.ReadFile(n.thermalPath + "/" + sensor + "/type")
    perr(err)

    temp, err := os.ReadFile(n.thermalPath + "/" + sensor + "/temp")
    if err != nil {
      continue
    }
    if string(temp[:len(temp)-1]) == "0" {
      continue
    }

    sensors = append(sensors, sensor)
  }

  if len(sensors) == 0 {
    perr(errThermals)
  }

  for i := 0; i < len(sensors); i++ {
    t.Register(NewNodeThermalSensor(n.thermalPath, sensors[i]))
  }
}

func NewNodeThermal(thermalPath string) *NodeThermal {
  n := new(NodeThermal)
  n.NodeBase = NewNodeBase()
  n.thermalPath = thermalPath
  return n
}

func (n *NodeThermal) Name() string {
  return "thermal"
}
