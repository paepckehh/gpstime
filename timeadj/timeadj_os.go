//go:build !freebsd && !linux && !darwin

package timeadj

import "syscall"

func setUnixTS(unixsec, rest int64) error {
	time := syscall.Timeval{
		Sec:  unixsec,
		Usec: rest,
	}
	err := syscall.Settimeofday(&time)
	return err
}
