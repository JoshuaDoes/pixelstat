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
	return fmt.Sprintf("%."+i64tostr(precision)+"f", in)
}

func strtof32(in string) float32 {
	num, err := strconv.ParseFloat(in, 32)
	perr(err)
	return float32(num)
}
func f32tostr(in float32, precision int64) string {
	return fmt.Sprintf("%."+i64tostr(precision)+"f", in)
}

func perr(err error) {
	if err != nil {
		fmt.Printf("\n\nERROR!\n\n%v\n\n", err)
		os.Exit(1)
	}
}

// Hertz converts a duration into a frequency rate.
func Hertz(rate time.Duration) int {
	return int(time.Second / rate)
}

// HertzPrecise converts a duration into a frequency rate with the highest precision that Go can muster.
func HertzPrecise(rate time.Duration) float64 {
	return float64(time.Second / rate)
}

// Rate converts a frequency in hertz to a time duration.
func Rate(hz int) time.Duration {
	return time.Second / time.Duration(hz)
}

// RatePrecise converts a frequency in hertz to a time duration with the highest precision that Go can muster.
func RatePrecise(hz float64) time.Duration {
	return time.Second / time.Duration(hz)
}

func Loop(sleep *time.Duration, fnc func()) (stopper chan bool, timeStart *time.Time) {
	stopper = make(chan bool)
	now := time.Now()
	loopThread(stopper, sleep, &now, fnc)
	return stopper, &now
}

func loopThread(stopper chan bool, sleep *time.Duration, timeStart *time.Time, fnc func()) {
	//It's expected for timeStart to fall behind if pacing is too fast,
	//could be exposed later on to report how far behind it falls.
	go func(stopper chan bool, sleep *time.Duration, timeStart *time.Time, fnc func()) {
		for {
			go fnc()
			time.Sleep(*sleep - time.Since(*timeStart))
			*timeStart = timeStart.Add(*sleep)
			select {
			case <-stopper:
				close(stopper)
				return
			default:
			}
		}
	}(stopper, sleep, timeStart, fnc)
}
