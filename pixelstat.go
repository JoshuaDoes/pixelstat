package main

import (
  "github.com/spf13/pflag"

  "os"
  "os/signal"
  "syscall"
//  "time"
)

var (
  pollingRate int  = 1000
  refreshRate int  = 240
  novrr       bool = false
  nothermal   bool = false
  netbits     bool = false
)

func main() {
  pflag.IntVar(&pollingRate, "poll", pollingRate, "max polling rate per node")
  pflag.IntVar(&refreshRate, "rate", refreshRate, "max refresh rate per renderer")
  pflag.BoolVar(&novrr, "novrr", false, "disable variable refresh rate")
  pflag.BoolVar(&nothermal, "nothermal", false, "disable thermal nodes")
  pflag.BoolVar(&netbits, "netbits", false, "use bits instead of bytes for netspeed")
  pflag.Parse()

  t := NewTracker()

  battery, err := NewNodeBattery("/sys/class/power_supply/battery/capacity")
  if err == nil { t.Register(battery) }
  charging, err := NewNodeCharging("/sys/class/power_supply/battery/status")
  if err == nil { t.Register(charging) }
  t.Register(NewNodeNetSpeed(netbits, "/sys/class/net",
    "wlan0", "wlan1", "rmnet2"))
  t.Register(NewNodeCPUFreq("/sys/devices/system/cpu/cpufreq"))
  t.Register(NewNodeCPUSet("/dev/cpuset"))
  voltage, err := NewNodeVoltage("/sys/class/power_supply/battery/voltage_now")
  if err == nil { t.Register(voltage) }
  current, err := NewNodeCurrent("/sys/class/power_supply/battery/current_now")
  if err == nil { t.Register(current) }
  t.Register(NewNodeWattage())

  if !nothermal {
    t.Register(NewNodeThermal("/sys/class/thermal",
      "LITTLE", "MID", "BIG",
      "G3D", "TPU", "soc",
      "gnss_tcxo_therm", "disp_therm",
      "usb_pwr_therm", "usb_pwr_therm2",
      "qi_therm", "battery", "batt_vs", "maxfg",
      "neutral_therm", "quiet_therm"))
  }

  t.Register(NewNodeRenderer(refreshRate, !novrr))
  t.Start(pollingRate)

  sig := make(chan os.Signal, 1)
  signal.Notify(sig, syscall.SIGINT)  //Keyboard interrupt
  signal.Notify(sig, syscall.SIGHUP)  //Terminal disappeared
  signal.Notify(sig, syscall.SIGKILL) //Process abandoned by kernel, how are we here???
  <-sig

  t.Close()
}
