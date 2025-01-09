package main

import (
  "github.com/JoshuaDoes/crunchio"
)

type Node interface {
  SetTracker(*Tracker)
  GetTracker() *Tracker
  Name()       string
  Rate()       int
  Unit()       string
  ValueType()  string
  ValueLen()   int64
  Value()      *crunchio.Buffer
  Close()      error
}
