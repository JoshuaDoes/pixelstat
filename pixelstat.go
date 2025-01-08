package main

import (
  "github.com/spf13/pflag"

  "os"
  "os/signal"
  "syscall"
)

var (
  refreshRate int  = 1920
  vrr         bool = true
)

func main() {
  pflag.IntVar(&refreshRate, "rate", refreshRate, "max refresh rate")
  pflag.BoolVar(&vrr, "vrr", vrr, "variable refresh rate")
  pflag.Parse()

  t := NewTracker(
    NewNodeBattery("/sys/class/power_supply/battery/capacity"),
    NewNodeCharging("/sys/class/power_supply/battery/status"),
    NewNodeVoltage("/sys/class/power_supply/battery/voltage_now"),
    NewNodeCurrent("/sys/class/power_supply/battery/current_now"),
    NewNodeWattage(),
    NewNodeCPUFreq("/sys/devices/system/cpu/cpufreq"),
    NewNodeThermal("/sys/class/thermal"),
    NewNodeRenderer(refreshRate, vrr),
  )
  t.Start()

  sig := make(chan os.Signal, 1)
  signal.Notify(sig, syscall.SIGINT)  //Keyboard interrupt
  signal.Notify(sig, syscall.SIGHUP)  //Terminal disappeared
  signal.Notify(sig, syscall.SIGKILL) //Process abandoned by kernel, how are we here???
  <-sig

  t.Close()
}
