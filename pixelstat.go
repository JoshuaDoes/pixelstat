package main

import (
  "github.com/spf13/pflag"

  "fmt"
  "os"
  "os/signal"
  "runtime"
  "syscall"
  "time"
)

var (
  pollingRate int  = 1000
  refreshRate int  = 120
  novrr       bool = false
  nocpuset    bool = false
  nonet       bool = false
  nothermal   bool = false
  netbits     bool = false
)

func main() {
  //defer recovery()

  pflag.IntVar(&pollingRate, "poll", pollingRate, "max polling rate per node")
  pflag.IntVar(&refreshRate, "rate", refreshRate, "max refresh rate per renderer")
  pflag.BoolVar(&novrr, "novrr", false, "disable variable refresh rate")
  pflag.BoolVar(&nocpuset, "nocpuset", false, "disable cpuset nodes")
  pflag.BoolVar(&nonet, "nonet", false, "disable networking nodes")
  pflag.BoolVar(&nothermal, "nothermal", false, "disable thermal nodes")
  pflag.BoolVar(&netbits, "netbits", false, "use bits instead of bytes for netspeed")
  pflag.Parse()

  t := NewTracker()

  batteryPaths := []string{"battery", "BAT0"}
  for _, path := range batteryPaths {
    battery, err := NewNodeBattery("/sys/class/power_supply/" + path + "/capacity")
    if err == nil {
      t.Register(battery)
      break
    }
  }
  for _, path := range batteryPaths {
    charging, err := NewNodeCharging("/sys/class/power_supply/" + path + "/status")
    if err == nil {
      t.Register(charging)
      break
    }
  }
  for _, path := range batteryPaths {
    path = "/sys/class/power_supply/" + path
    voltage, err := NewNodeVoltage(path + "/voltage_now")
    if err != nil {
      continue
    }
    current, err := NewNodeCurrent(path + "/current_now")
    if err != nil {
      continue
    }
    t.Register(voltage, current, NewNodeWattage())
    break
  }

  if gpufreq, err := NewNodeGPUFreq("/sys/class/devfreq/34f00000.gpu0/cur_freq", 1000000); err == nil {
    t.Register(gpufreq)
  } else if gpufreq, err := NewNodeGPUFreq("/sys/devices/platform/1c500000.mali/cur_freq", 1000); err == nil {
    t.Register(gpufreq)
  }

  t.Register(NewNodeCPUFreq("/sys/devices/system/cpu/cpufreq"))

  t.Register(NewNodeCPUBandwidth())

  if !nocpuset {
    t.Register(NewNodeCPUSet("/dev/cpuset"))
  }

  if !nonet {
    t.Register(NewNodeNetSpeed(netbits, "/sys/class/net"))
  }

  if !nothermal {
    t.Register(NewNodeThermal("/sys/class/thermal"))
  }

  renderer := NewNodeRenderer(refreshRate, !novrr)
  t.Register(renderer)

  t.Start(pollingRate)
  sig := make(chan os.Signal, 1)
  signal.Notify(sig, syscall.SIGINT)  //Keyboard interrupt
  signal.Notify(sig, syscall.SIGHUP)  //Terminal disappeared
  signal.Notify(sig, syscall.SIGKILL) //Process abandoned by kernel, how are we here???
  <-sig

  go t.Close()
  fmt.Printf("\nShutting down...\n")
  for {
    time.Sleep(time.Millisecond * 1)
    if t.Closed() {
      break
    }
  }
  fmt.Println("Goodbye!")
}

func recovery() {
  if reason := recover(); reason != nil {
    stack := make([]byte, 65536)
    stackfull := make([]byte, 65536)
    stackN := runtime.Stack(stack, false)
    stackfullN := runtime.Stack(stackfull, true)
    os.WriteFile("stack.log", stack[:stackN], 0644)
    os.WriteFile("stackfull.log", stackfull[:stackfullN], 0644)
    os.Stderr.Write(stack)
    os.Exit(1)
  }
}
