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
  if err != nil {
    return
  }

  sensors := make([]string, 0)
  for i := 0; i < len(dir); i++ {
    file := dir[i]
    if file.Type() != fs.ModeSymlink {
      perr(errThermalsNotSymlink)
    }

    sensor := file.Name()
    found := false

    name, err := os.ReadFile(n.thermalPath + "/" + sensor + "/type")
    perr(err)
    nameS := string(name[:len(name)-1])
    for j := 0; j < len(n.sensors); j++ {
      if n.sensors[j] == nameS {
        found = true
        break
      }
    }
    if !found {
      continue
    }

    _, err = os.ReadFile(n.thermalPath + "/" + sensor + "/temp")
    if err != nil {
      continue
    }

    sensors = append(sensors, sensor)
  }

  for i := 0; i < len(sensors); i++ {
    node, err := NewNodeThermalSensor(n.thermalPath, sensors[i])
    if err == nil {
      t.Register(node)
    }
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
