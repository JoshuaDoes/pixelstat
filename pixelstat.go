package main

import (
  "os"
  "os/signal"
  "syscall"
)

func main() {
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
  stopper := r.RefreshTerminal(120)
  _ = stopper

  sig := make(chan os.Signal, 1)
  signal.Notify(sig, syscall.SIGINT)  //Keyboard interrupt
  signal.Notify(sig, syscall.SIGHUP)  //Terminal disappeared
  signal.Notify(sig, syscall.SIGKILL) //Process abandoned by kernel, how are we here???
  <-sig

  //stopper <- true
  //t.Close()
}
