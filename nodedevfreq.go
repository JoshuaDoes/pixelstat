package main

import (
	"os"
)

type NodeDevFreq struct {
	*NodeBase
	devFreqPath string
}

func NewNodeDevFreq(devFreqPath string) *NodeDevFreq {
	n := new(NodeDevFreq)
	n.NodeBase = NewNodeBase()
	n.devFreqPath = devFreqPath
	return n
}

func (n *NodeDevFreq) Name() string {
	return "devfreq"
}

func (n *NodeDevFreq) SetTracker(t *Tracker) {
	dir, err := os.ReadDir(n.devFreqPath)
	if err != nil {
		return
	}

	devs := make([]string, 0)
	for _, file := range dir {
		dev := file.Name()
		freqFile, err := os.Open(n.devFreqPath + "/" + dev + "/cur_freq")
		if err != nil {
			continue
		}
		freqFile.Close()
		devs = append(devs, dev)
	}

	for i := 0; i < len(devs); i++ {
		if node, err := NewNodeDevFreqDevice(n.devFreqPath, devs[i]); err == nil {
			t.Register(node)
		}
	}
}
