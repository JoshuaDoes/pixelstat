package main

import (
	"github.com/JoshuaDoes/crunchio"

	"sync"
	"time"
)

type TrackerValue struct {
	sync.Mutex
	stopper chan bool
	node    Node

	value   *crunchio.Buffer
	cnt     int64
	sum     []float64
	min     []float64
	max     []float64
	indexes int

	ppsTime time.Time
	ppsLast float64
	pps     int64
}

func (tv *TrackerValue) start(maxHz int) {
	defer recovery()

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
		hertz := Rate(rate)
		tv.stopper, _ = Loop(&hertz, tv.getValue)
	}
}

func (tv *TrackerValue) close() {
	perr(tv.node.Close())
}

func (tv *TrackerValue) getValue() {
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
				tv.ppsLast = float64(tv.pps) / (float64(delta.Nanoseconds()) / 1000000000)
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
				tv.expandIndexes(len(i32s))
				for i := 0; i < len(i32s); i++ {
					tv.setValue(i, float64(i32s[i]))
				}
			case "i64":
				length := n.ValueLen()
				if bc := v.ByteCapacity(); (bc / 8) < length {
					length = (bc / 8)
				}
				i64s := v.ReadI64LENext(length)
				tv.expandIndexes(len(i64s))
				for i := 0; i < len(i64s); i++ {
					tv.setValue(i, float64(i64s[i]))
				}
			case "f32":
				length := n.ValueLen()
				if bc := v.ByteCapacity(); (bc / 4) < length {
					length = (bc / 4)
				}
				f32s := v.ReadF32LENext(length)
				tv.expandIndexes(len(f32s))
				for i := 0; i < len(f32s); i++ {
					tv.setValue(i, float64(f32s[i]))
				}
			case "f64":
				length := n.ValueLen()
				if bc := v.ByteCapacity(); (bc / 8) < length {
					length = (bc / 8)
				}
				f64s := v.ReadF64LENext(length)
				tv.expandIndexes(len(f64s))
				for i := 0; i < len(f64s); i++ {
					tv.setValue(i, f64s[i])
				}
			}
		}(crunchio.NewBuffer(v.Bytes()), tv)
	}
}

func (tv *TrackerValue) setValue(i int, value float64) {
	tv.sum[i] += value
	if tv.min[i] == 0 || tv.min[i] > value {
		tv.min[i] = value
	}
	if tv.max[i] == 0 || tv.max[i] < value {
		tv.max[i] = value
	}
}

func (tv *TrackerValue) Value() *crunchio.Buffer {
	if tv == nil {
		return nil
	}
	if tv.value == nil {
		tv.getValue()
	}
	if tv.value == nil {
		return nil
	}

	return crunchio.NewBuffer(tv.value.Bytes())
}

func (tv *TrackerValue) Average(index int64) float64 {
	if tv.cnt == 0 || len(tv.sum) == 0 {
		return 0
	}
	return tv.sum[index] / float64(tv.cnt)
}

func (tv *TrackerValue) Min(index int64) float64 {
	if tv.cnt == 0 || len(tv.min) == 0 {
		return 0
	}
	return tv.min[index]
}

func (tv *TrackerValue) Max(index int64) float64 {
	if tv.cnt == 0 || len(tv.max) == 0 {
		return 0
	}
	return tv.max[index]
}

func (tv *TrackerValue) Polls() int64 {
	return tv.cnt
}

func (tv *TrackerValue) PollsPerSecond() float64 {
	return tv.ppsLast
}

func (tv *TrackerValue) expandIndexes(length int) {
	if tv.indexes == 0 {
		tv.sum = make([]float64, length)
		tv.min = make([]float64, length)
		tv.max = make([]float64, length)
	} else if length > tv.indexes {
		tv.sum = append(tv.sum, make([]float64, length-tv.indexes)...)
		tv.min = append(tv.min, make([]float64, length-tv.indexes)...)
		tv.sum = append(tv.sum, make([]float64, length-tv.indexes)...)
	}
	tv.indexes = length
}
