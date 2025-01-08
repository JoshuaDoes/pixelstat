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
  sensors []string
}

func (n *NodeThermal) SetTracker(t *Tracker) {
  dir, err := os.ReadDir(n.thermalPath)
  perr(err)

  sensors := make([]string, 0)
  for i := 0; i < len(dir); i++ {
    file := dir[i]
    if file.Type() != fs.ModeSymlink {
      perr(errThermalsNotSymlink)
    }

    sensor := file.Name()
    found := false
    for j := 0; j < len(n.sensors); j++ {
      if n.sensors[j] == sensor {
        found = true
        break
      }
    }
    if !found {
      continue
    }

    _, err = os.ReadFile(n.thermalPath + "/" + sensor + "/type")
    perr(err)

    _, err = os.ReadFile(n.thermalPath + "/" + sensor + "/temp")
    if err != nil {
      continue
    }

    sensors = append(sensors, sensor)
  }

  for i := 0; i < len(sensors); i++ {
    t.Register(NewNodeThermalSensor(n.thermalPath, sensors[i]))
  }
}

func NewNodeThermal(thermalPath string, sensors ...string) *NodeThermal {
  n := new(NodeThermal)
  n.NodeBase = NewNodeBase()
  n.thermalPath = thermalPath
  n.sensors = sensors
  return n
}

func (n *NodeThermal) Name() string {
  return "thermal"
}
