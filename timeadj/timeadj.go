package timeadj

import "time"

const (
	_lf      = "\n"
	_ts      = "2006-01-02 15:04:05.999"
	_minTime = 1662897652
)

// SecureTestAdjust ...
func SecureTestAdjust(tsTarget time.Time) (bool, string) {
	tsLocal := time.Now()
	diff := tsTarget.Sub(tsLocal)
	diffSeconds := diff.Seconds()
	if diffSeconds > 0.8 || diffSeconds < -0.8 { // allow +/- 800ms [rounding|jitter], bail out [early|cheap|eco-friendly]
		if tsTarget.IsZero() {
			return false, "zero"
		}
		if tsTarget.Unix() < _minTime {
			return false, "error target time before minimum time [buildTime]"
		}
		if diffSeconds > 0 {
			if err := setTS(tsTarget); err != nil {
				return false, "[error] [" + err.Error() + "]"
			}
			if diffSeconds > 30 {
				msg := "[time] [adjust] [forward] [large] [init] [diff: +" + tsTarget.Sub(tsLocal).String() + "]" + _lf
				msg += "old  time: " + tsLocal.Format(_ts) + _lf
				msg += "new  time: " + tsTarget.Format(_ts) + _lf
				return true, msg
			}
			return true, "[time] [adjust] [forward] [diff: +" + diff.String() + "]" + _lf
		}
		if diffSeconds < -10 {
			msg := "[time] [adjust] [backward] [large] [ALERT] [DENY] [diff: " + diff.String() + "]" + _lf
			msg += "old : " + tsLocal.Format(_ts) + _lf
			msg += "new : " + tsTarget.Format(_ts) + _lf
			return false, msg
		}
		if err := setTS(tsTarget); err != nil {
			return false, "[error] [" + err.Error() + "]"
		}
		return true, "[time] [adjust] [backward] [diff: " + diff.String() + "]" + _lf
	}
	return false, ""
}

func setTS(ts time.Time) error {
	sec := ts.Unix()
	return setUnixTS(sec, ts.UnixMicro()-(sec*1000*1000))
}
