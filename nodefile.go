package main

import (
  "github.com/JoshuaDoes/crunchio"

  "os"
  "sync"
)

type NodeFile struct {
  sync.Mutex
  *NodeBase
  Path string
  File *os.File
}

func NewNodeFile(path string) *NodeFile {
  n := new(NodeFile)
  n.NodeBase = NewNodeBase()
  n.Path = path

  f, err := os.Open(path)
  perr(err)
  n.File = f

  return n
}

func (n *NodeFile) Name() string {
  return "file"
}

func (n *NodeFile) Value() *crunchio.Buffer {
  n.Lock()
  defer n.Unlock()
  //data, err := os.ReadFile(n.Path)
  //perr(err)

  buf := make([]byte, 32768)
  data := make([]byte, 0)

  for {
    read, err := n.File.Read(buf)
    if err != nil {
      n.File.Seek(0, 0)
      break
    }
    data = append(data, buf[:read]...)
  }

  return crunchio.NewBuffer(data)
}
