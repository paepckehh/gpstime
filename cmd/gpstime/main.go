package main

import (
	"os"

	"paepcke.de/gpstime"
	"paepcke.de/gpstime/gpsfeed"
)

const (
	_defaultMode   = "time"
	_defaultDevice = "/dev/gps0"
)

// gpstime
func main() {
	mode, port := _defaultMode, _defaultDevice
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "time":
		case "location":
			mode = "location"
		default:
			port = gpsfeed.GetDeviceName(os.Args[i], port)
		}
	}
	switch mode {
	case "time":
		gpstime.CheckSetClock(port)
	case "location":
		gpstime.CheckSetClockLocation(port)
	}
	out("[error] [internal error - undefined mode] [please report]")
}

//
// Little Helper
//

// out ...
func out(msg string) { os.Stdout.Write([]byte(msg)) }
