package main

import (
  "github.com/spf13/pflag"

  "os"
  "os/signal"
  "syscall"
)

var (
  pollingRate int  = 1920
  refreshRate int  = 1920
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

  t := NewTracker(
    NewNodeBattery("/sys/class/power_supply/battery/capacity"),
    NewNodeCharging("/sys/class/power_supply/battery/status"),
    NewNodeNetSpeed(netbits, "/sys/class/net",
      "wlan0", "wlan1", "rmnet2"),
    NewNodeCPUFreq("/sys/devices/system/cpu/cpufreq"),
    NewNodeVoltage("/sys/class/power_supply/battery/voltage_now"),
    NewNodeCurrent("/sys/class/power_supply/battery/current_now"),
    NewNodeWattage(),
  )

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
