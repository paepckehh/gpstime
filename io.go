package gpstime

import (
	"os"
)

// const
const (
	_inf    = "[gpstime] [info] "
	_err    = "[gpstime] [error] "
	_errMax = 10 // max number of device error reports [resets at success]
)

//
// DISPLAY IO
//

// out handles output messages to stdout, adding an linefeed
func out(msg string) { outPlain(msg + "\n") }

// outPlain handles output messages to stdout
func outPlain(msg string) { os.Stdout.Write([]byte(msg)) }
