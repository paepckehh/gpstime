//go:build darwin

package timeadj

import "syscall"

func setUnixTS(unixsec, rest int64) error {
	time := syscall.Timeval{
		Sec:  unixsec,
		Usec: int32(rest),
	}
	err := syscall.Settimeofday(&time)
	return err
}
