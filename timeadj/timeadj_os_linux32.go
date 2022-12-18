//go:build (linux && arm) || (linux && i386)

package timeadj

import "syscall"

func setUnixTS(unixsec, rest int64) error {
	time := syscall.Timeval{
		Sec:  int32(unixsec),
		Usec: int32(rest),
	}
	err := syscall.Settimeofday(&time)
	return err
}
