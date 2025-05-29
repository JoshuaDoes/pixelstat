package main

import (
  "fmt"
  "os"
  "strconv"
  "time"
)

func strtoi64(in string) int64 {
  num, err := strconv.ParseInt(in, 10, 64)
  perr(err)
  return num
}
func i64tostr(in int64) string {
  return fmt.Sprintf("%d", in)
}

func strtoi32(in string) int32 {
  num, err := strconv.ParseInt(in, 10, 32)
  perr(err)
  return int32(num)
}
func i32tostr(in int32) string {
  return fmt.Sprintf("%d", in)
}

func strtof64(in string) float64 {
  num, err := strconv.ParseFloat(in, 64)
  perr(err)
  return num
}
func f64tostr(in float64, precision int64) string {
  return fmt.Sprintf("%." + i64tostr(precision) + "f", in)
}

func strtof32(in string) float32 {
  num, err := strconv.ParseFloat(in, 32)
  perr(err)
  return float32(num)
}
func f32tostr(in float32, precision int64) string {
  return fmt.Sprintf("%." + i64tostr(precision) + "f", in)
}

func perr(err error) {
  if err != nil {
    fmt.Printf("\n\nERROR!\n\n%v\n\n", err)
    os.Exit(1)
  }
}

func hertz(hz int) time.Duration {
  return time.Second / time.Duration(hz)
}

func loop(rate int, fnc func()) chan bool {
  defer recovery()

  stopper := make(chan bool)
  cancel  := make(chan bool)

  go loopSignal(stopper, cancel)
  go loopRunner(cancel, rate, fnc)

  return stopper
}

func loopSignal(stopper, cancel chan bool) {
  defer recovery()
  stop := <- stopper
  cancel <- stop
}

func loopRunner(cancel chan bool, rate int, fnc func()) {
  defer recovery()
  pace := hertz(rate)

  //It's expected for timeStart to fall behind if pacing is too fast,
  //could be exposed later on to report how far behind it falls.
  timeStart := time.Now()
  for {
    fnc()
    time.Sleep(pace - time.Since(timeStart))
    timeStart = timeStart.Add(pace)

    select {
    case <- cancel:
      close(cancel)
      return
    default:
      //pass
    }
  }
}
