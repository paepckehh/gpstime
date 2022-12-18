package gpstime

import (
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"paepcke.de/airloctag"
	"paepcke.de/daylight/sun"
	"paepcke.de/gpstime/gpsfeed"
)

const (
	_TS          = "15:04:05" // time stamp layout [time.Parse]
	_targetDir   = "/var/gps"
	_targetFile  = "/var/gps/.location"
	_header      = "#!/bin/sh\nexport GPS_MODE=\"gpstime\"\n"
	_lat         = "export GPS_LAT=\""
	_long        = "export GPS_LONG=\""
	_elev        = "export GPS_ELEVATION=\""
	_atag        = "export AIRLOCTAG=\""
	_sunrise     = "export GPS_SUN_RISE=\""
	_sunset      = "export GPS_SUN_SET=\""
	_noon        = "export GPS_SUN_NOON=\""
	_daylight    = "export GPS_SUN_DAYLIGHT=\""
	_optLongest  = "export GPS_SUN_OPT=\"[-=* LONGEST DAY OF THE YEAR *=-]"
	_optShortest = "export GPS_SUN_OPT=\"[-=* SHORTEST DAY OF THE YEAR *=-]"
	_report      = "[location] [adjust]"
	_lf          = "\n"       // linefeed
	_qlf         = "\"" + _lf // escaped doubleqoute [end] and linefeed
)

// writeLocaltionFile reads gps nmea frames from a gps device and writes the current location as unix shell env include
func writeLocationFile(dev *gpsfeed.GpsDevice) {
	delay, initDone := gpsfeed.DeviceTimeout, false
	for {
		if writeEnvFile(parseLocation(dev)) {
			dev.ErrCount.Store(0)
			initDone = true
			delay = gpsfeed.DeviceTimeout
		}
		if delay < 1440 {
			if initDone {
				delay++
			}
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
	dev.Global.Done() // currently unreachable
}

// writeLocaltionFileOnce exits after success
func writeLocationFileOnce(dev *gpsfeed.GpsDevice) {
	for {
		if writeEnvFile(parseLocation(dev)) {
			break
		}
		time.Sleep(time.Duration(gpsfeed.DeviceTimeout) * time.Second)
	}
}

// parseLocation parses gps devices sentences, extracts a location [lat,long,elev]
func parseLocation(dev *gpsfeed.GpsDevice) (c Coord) {
	dev.Open()
	defer dev.Close()
	for dev.Feed.Scan() {
		line := dev.Feed.Text()
		dev.Responsive.Store(true)
		if len(line) > 12 {
			if line[:6] == "$GPGGA" {
				dev.DataValid.Store(true) // dearm watchdog
				if c = parseLocationStamp(line); c.Valid {
					break // success, return
				}
				if dev.CheckErrAdd() {
					break // fail, return
				}
			}
		}
	}
	return c
}

// parseLocationStamp parses a single gps GGA sentence
func parseLocationStamp(line string) (c Coord) {
	var ok bool
	var err error
	s := strings.Split(line, ",")
	if len(s) != 15 {
		out(_err + " [invalid token] [nmea.GGA] [" + line + "]")
		return c
	}
	if ok, c.Lat = parseGPS(s[2], s[3]); !ok {
		out(_err + " [unable to parse] [nmea.GGA] [lat] [" + line + "]")
		return c
	}
	if ok, c.Long = parseGPS(s[4], s[5]); !ok {
		out(_err + " [unable to parse] [nmea.GGA] [long] [" + line + "]")
		return c
	}
	if c.Elevation, err = strconv.ParseFloat(s[9], 64); err != nil {
		out(_err + " [unable to parse] [nmea.GGA] [elevation] [" + line + "]")
		return c
	}
	c.Valid = true
	return c
}

// parseGPS validates and parsrs a degree & orientation GPS pair
func parseGPS(val, orientation string) (bool, float64) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false, 0
	}
	d := math.Floor(v / 100)
	m := v - (d * 100)
	switch orientation {
	case "N", "E":
		return true, d + m/60
	case "S", "W":
		return true, 0 - (d + m/60)
	}
	return false, 0
}

var oldPoint Coord

func fl(in float64) string { return strconv.FormatFloat(in, 'f', 8, 64) }
func fs(in float64) string { return strconv.FormatFloat(in, 'f', 1, 64) }

// writeEnvFile takes a location [lat,long,elevation] and writes it as a unix shell variable include
func writeEnvFile(c Coord) bool {
	if !pointChanged(c) {
		return false
	}
	_ = os.MkdirAll(_targetDir, 0o444)
	err := os.WriteFile(_targetFile, buildEnv(c), 0o444)
	if err != nil {
		out(_err + "[unable to write gps env file] [" + err.Error() + "}")
	}
	out(_inf + _report + " [lat: " + fl(c.Lat) + "] [long: " + fl(c.Long) + "] [elevation: " + fl(c.Elevation) + "]")
	return true
}

// builds the env script include
func buildEnv(c Coord) []byte {
	sunrise, sunset, noon, daylight, longest, shortest := sun.StateExtended(c.Lat, c.Long, c.Elevation)
	hash, _, _ := airloctag.Encode(c.Lat, c.Long, c.Elevation, "", 0)
	return []byte(_header + _lat + fl(c.Lat) + _qlf + _long + fl(c.Long) + _qlf + _elev + fs(c.Elevation) + _qlf + _atag + hash + _qlf + _sunrise + sunrise.Format(_TS) + _qlf + _sunset + sunset.Format(_TS) + _qlf + _noon + noon.Format(_TS) + _qlf + _daylight + daylight.String() + _qlf + getOpt(longest, shortest) + _lf)
}

func r1(in float64) float64 { return math.Round(in*10) / 10 }
func r2(in float64) float64 { return math.Round(in*1000) / 1000 }

// roundCoord to remove jitter
func roundCoord(in Coord) Coord {
	return Coord{false, r2(in.Lat), r2(in.Long), r1(in.Elevation)}
}

// pointChanged checks if we moved
func pointChanged(c Coord) bool {
	round := roundCoord(c)
	if oldPoint.Lat == round.Lat && oldPoint.Long == round.Long {
		return false
	}
	oldPoint.Lat, oldPoint.Long, oldPoint.Elevation = round.Lat, round.Long, round.Elevation // set new global reference
	return true
}

// getOpt translate sun turning point indicator into env var strings
func getOpt(longest, shortest bool) string {
	switch {
	case longest:
		return _optLongest
	case shortest:
		return _optShortest
	}
	return ""
}
