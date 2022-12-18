//go:build (freebsd && arm) || (freebsd && i386)

package timeadj

import (
	"errors"
	"syscall"
)

func setUnixTS(unixsec, rest int64) error {
	s, err := syscall.Sysctl("kern.securelevel")
	if err != nil {
		return errors.New("unable to verify kern.securelevel [" + err.Error() + "]")
	}
	if s == "2" || s == "3" {
		return errors.New("kern.securelevel " + s + " detected! [ ->largest timejump 1sec]")
	}
	time := syscall.Timeval{
		Sec:  unixsec,
		Usec: int32(rest),
	}
	err = syscall.Settimeofday(&time)
	return err
}
