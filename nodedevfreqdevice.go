package main

import (
  "github.com/JoshuaDoes/crunchio"

  "sync"
)

type NodeDevFreqDevice struct {
	sync.Mutex
	*NodeFile

	device string
}

func NewNodeDevFreqDevice(devFreqPath, device string) (*NodeDevFreqDevice, error) {
	n := new(NodeDevFreqDevice)
	n.device = device
	path := devFreqPath + "/" + device
	nf, err := NewNodeFile(path + "/cur_freq")
	n.NodeFile = nf
	return n, err
}

func (n *NodeDevFreqDevice) Name() string {
	return n.device + " freq"
}

func (n *NodeDevFreqDevice) Rate() int {
	return 1000
}

func (n *NodeDevFreqDevice) Unit() string {
	return " MHz"
}

func (n *NodeDevFreqDevice) ValueType() string {
	return "i64"
}

func (n *NodeDevFreqDevice) ValueLen() int64 {
	return 1
}

func (n *NodeDevFreqDevice) Value() *crunchio.Buffer {
	n.Lock()
	defer n.Unlock()

  v := n.NodeFile.Value()
  if v == nil || v.ByteCapacity() == 0 {
    return nil
  }
  v.TruncateRight(1) //Remove the newline

  MHz := strtoi64(v.String()) / 1000 / 1000

  v = crunchio.NewBuffer()
  v.Grow(8)
  v.WriteAbstract(MHz)
  v.Seek(0, 0)
  return v
}
