package main

import (
  "github.com/JoshuaDoes/crunchio"

  "time"
)

type Node interface {
  SetTracker(*Tracker)
  GetTracker() *Tracker
  Name()      string
  Rate()      time.Duration
  Unit()      string
  ValueLen()  int
  Value()     *crunchio.Buffer
}
