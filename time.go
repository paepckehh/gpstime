package gpstime

import (
	"strings"
	"time"

	"paepcke.de/gpstime/gpsfeed"
	"paepcke.de/gpstime/timeadj"
)

const (
	_ts    = "20060102-150405.999" // time date stamp layout [time.Parse]
	_empty = ""
	_zero  = "zero"
	_error = "error"
)

// CheckSetClock reads raw GPS NEMEA $GP frames from device and adjusts the local system time
func checkSetClock(dev *gpsfeed.GpsDevice) {
	delay, timeChanged, msg := gpsfeed.DeviceTimeout, false, ""
	for {
		if timeChanged, msg = timeadj.SecureTestAdjust(parseTime(dev)); timeChanged {
			delay = gpsfeed.DeviceTimeout
		}
		if msg != _empty {
			if len(msg) > 3 {
				if msg == _zero {
					continue
				}
			}
			if len(msg) > 4 {
				if msg[:5] == _error {
					dev.ErrCount.Add(1)
					out(_err + msg)
					continue
				}
			}
			outPlain(_inf + msg)
		}
		if !dev.InitDone.Load() {
			dev.InitDone.Store(true)
		}
		if delay < 1440 {
			delay++
		}
		dev.ErrCount.Store(0)
		time.Sleep(time.Duration(delay) * time.Second)
	}
	// unreachable dev.Global.Done()
}

// CheckSetClockOnce reads raw GPS NEMEA $GP sentences from device and adjusts the local system time and exits when successful
func checkSetClockOnce(dev *gpsfeed.GpsDevice) {
	for {
		_, msg := timeadj.SecureTestAdjust(parseTime(dev))
		if msg != _empty {
			if len(msg) > 3 {
				if msg == _zero {
					time.Sleep(time.Duration(gpsfeed.DeviceTimeout) * time.Second)
					continue
				}
			}
			if len(msg) > 4 {
				if msg[:5] == _error {
					dev.ErrCount.Add(1)
					out(_err + msg)
					time.Sleep(time.Duration(gpsfeed.DeviceTimeout) * time.Second)
					continue
				}
			}
			outPlain(_inf + msg)
		}
		out(_inf + "[time aleady in sync]")
		break
	}
}

// checkSetClockLocation
func checkSetClockLocation(dev *gpsfeed.GpsDevice) {
	dev.Global.Add(1) // exit if either one dies
	go checkSetClock(dev)
	for {
		if dev.InitDone.Load() {
			go writeLocationFile(dev)
			break
		}
		time.Sleep(1 * time.Second) // let timesync win this race for init
	}
	dev.Global.Wait()
}

// parseTimeFeed parses gps sentences -> time.Time stamps
func parseTime(dev *gpsfeed.GpsDevice) time.Time {
	dev.Open()
	defer dev.Close()
	for dev.Feed.Scan() {
		line := dev.Feed.Text()
		dev.Responsive.Store(true) // de-arm responsive watchdog
		if len(line) > 12 {
			if line[:6] == "$GPRMC" {
				dev.DataValid.Store(true) // de-arm data valid watchdog
				s := strings.Split(line, ",")
				switch {
				case len(s) != 13:
					if dev.CheckErrAdd() {
						out(_err + " [invalid nmea RMC token] [" + line + "]")
						continue
					}
				case len(s[9]) != 6:
					if dev.CheckErrAdd() {
						out(_err + "[invalid nmea RMC data stamp] [" + line + "]")
						continue
					}
				case len(s[1]) != 10:
					if dev.CheckErrAdd() {
						out(_err + " [invalid nmea RMC time stamp] [" + line + "]")
						continue
					}
				}
				stamp := "20" + s[9][4:] + s[9][2:4] + s[9][:2] + "-" + s[1]
				ts, err := time.Parse(_ts, stamp)
				if err != nil {
					if dev.CheckErrAdd() {
						out(_err + "[nmea RMC token] [time.Parse] [" + dev.GetErr() + "] [" + stamp + "]")
						continue
					}
				}
				return ts
			}
		}
	}
	return time.Time{} // faild to parse (return time.IsZero true)
}
