package main

import (
  "github.com/JoshuaDoes/crunchio"

  "fmt"
  "sync"
  "time"
)

type TrackerValue struct {
  sync.Mutex
  stopper chan bool
  node    Node

  value *crunchio.Buffer
  sum []float64
  cnt int64

  ppsTime time.Time
  ppsLast float64
  pps     int64
}

func (tv *TrackerValue) start(maxHz int) {
  if tv.stopper != nil {
    return
  }

  if tv.node.Rate() > 0 {
    rate := tv.node.Rate()
    if tv.node.ValueType() != "" {
      rateMax := maxHz
      if rate > rateMax {
        rate = rateMax
      }
    }
    tv.stopper = loop(rate, tv.getValue)
  }
}

func (tv *TrackerValue) close() {
  tv.stopper <- true
  tv.stopper = nil
  perr(tv.node.Close())
}

func (tv *TrackerValue) getValue() {
  defer func() {
    if r := recover(); r != nil {
      //Catch-all to ignore panics when reading nodes
      rs := fmt.Sprintf("%v", r)
      tv.value = crunchio.NewBuffer()
      tv.value.Grow(int64(len(rs)))
      tv.value.WriteAbstract(rs)
      tv.cnt++
    }
  }()

  tv.Lock()
  defer tv.Unlock()

  if tv.ppsTime.IsZero() {
    tv.ppsTime = time.Now()
  }

  n := tv.node
  v := n.Value()
  if v != nil && v.ByteCapacity() > 0 {
    tv.value = v
    tv.cnt++

    go func(tv *TrackerValue) {
      tv.Lock()
      defer tv.Unlock()

      tv.pps++
      delta := time.Since(tv.ppsTime)
      if delta.Seconds() >= 1 {
        tv.ppsTime = time.Now()
        tv.ppsLast = float64(tv.pps) / (((float64(delta.Nanoseconds()) / 1000) / 1000) / 1000)
        tv.pps = 0
      }
    }(tv)

    go func(v *crunchio.Buffer, tv *TrackerValue) {
      tv.Lock()
      defer tv.Unlock()

      switch tv.node.ValueType() {
      case "i32":
        length := n.ValueLen()
        if bc := v.ByteCapacity(); (bc / 4) < length {
          length = (bc / 4)
        }
        i32s := v.ReadI32LENext(length)
        tv.expandSum(len(i32s))
        for i := 0; i < len(i32s); i++ {
          tv.sum[i] += float64(i32s[i])
        }
      case "i64":
        length := n.ValueLen()
        if bc := v.ByteCapacity(); (bc / 8) < length {
          length = (bc / 8)
        }
        i64s := v.ReadI64LENext(length)
        tv.expandSum(len(i64s))
        for i := 0; i < len(i64s); i++ {
          tv.sum[i] += float64(i64s[i])
        }
      case "f32":
        length := n.ValueLen()
        if bc := v.ByteCapacity(); (bc / 4) < length {
          length = (bc / 4)
        }
        f32s := v.ReadF32LENext(length)
        tv.expandSum(len(f32s))
        for i := 0; i < len(f32s); i++ {
          tv.sum[i] += float64(f32s[i])
        }
      case "f64":
        length := n.ValueLen()
        if bc := v.ByteCapacity(); (bc / 8) < length {
          length = (bc / 8)
        }
        f64s := v.ReadF64LENext(length)
        tv.expandSum(len(f64s))
        for i := 0; i < len(f64s); i++ {
          tv.sum[i] += f64s[i]
        }
      }
    }(crunchio.NewBuffer(v.Bytes()), tv)
  }
}

func (tv *TrackerValue) Value() *crunchio.Buffer {
  if tv.value == nil {
    tv.getValue()
  }
  if tv.value == nil {
    return nil
  }

  tv.Lock()
  defer tv.Unlock()
  return crunchio.NewBuffer(tv.value.Bytes())
}

func (tv *TrackerValue) Average(index int64) float64 {
  if tv.cnt <= 0 || len(tv.sum) == 0 {
    return 0
  }
  return tv.sum[index] / float64(tv.cnt)
}

func (tv *TrackerValue) Polls() int64 {
  return tv.cnt
}

func (tv *TrackerValue) PollsPerSecond() float64 {
  return tv.ppsLast
}

func (tv *TrackerValue) expandSum(length int) {
  if tv.sum == nil {
    tv.sum = make([]float64, length)
    return
  }
  if length > len(tv.sum) {
    tv.sum = append(tv.sum, make([]float64, length - len(tv.sum))...)
  }
}
