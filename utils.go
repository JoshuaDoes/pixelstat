package main

import (
  "fmt"
  "os"
  "strconv"
  "time"
)

func hertz(hz time.Duration) time.Duration {
  return time.Second / hz
}

func str2int64(in string) int64 {
  num, err := strconv.ParseInt(in, 10, 64)
  perr(err)
  return num
}
func int642str(in int64) string {
  return fmt.Sprintf("%d", in)
}

func str2float(in string) float64 {
  num, err := strconv.ParseFloat(in, 64)
  perr(err)
  return num
}
func float2str(in float64, precision int64) string {
  return fmt.Sprintf("%." + int642str(precision) + "f", in)
}

func perr(err error) {
  if err != nil {
    fmt.Printf("\n\nERROR!\n\n%v\n\n", err)
    os.Exit(1)
  }
}

func loop(pace time.Duration, fnc func()) chan bool {
  stopper := make(chan bool)
  go func(stopper chan bool, pace time.Duration, fnc func()) {
    for {
      select {
      case <- stopper:
        return
      default:
        deadline := time.Now().Add(pace)
        fnc()
        time.Sleep(time.Until(deadline))
      }
    }
  }(stopper, pace, fnc)
  return stopper
}
