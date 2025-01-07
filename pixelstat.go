package main

import (
  "github.com/spf13/pflag"

  "os"
  "os/signal"
  "syscall"
)

var (
  refreshRate int = 480
)

func main() {
  pflag.IntVar(&refreshRate, "rate", refreshRate, "max refresh rate")
  pflag.Parse()

  r := NewRenderer()

  t := NewTracker(
    NewNodeBattery("/sys/class/power_supply/battery/capacity"),
    NewNodeCharging("/sys/class/power_supply/battery/status"),
    NewNodeVoltage("/sys/class/power_supply/battery/voltage_now"),
    NewNodeCurrent("/sys/class/power_supply/battery/current_now"),
    NewNodeWattage(),
    NewNodeCPUFreq("/sys/devices/system/cpu/cpufreq"),
    NewNodeThermal("/sys/class/thermal"),
    r,
  )

  t.Start()
  stopper := r.RefreshTerminal(refreshRate)
  _ = stopper

  sig := make(chan os.Signal, 1)
  signal.Notify(sig, syscall.SIGINT)  //Keyboard interrupt
  signal.Notify(sig, syscall.SIGHUP)  //Terminal disappeared
  signal.Notify(sig, syscall.SIGKILL) //Process abandoned by kernel, how are we here???
  <-sig

  //stopper <- true
  //t.Close()
}
