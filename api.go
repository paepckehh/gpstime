// package gpstime ...
package gpstime

import (
	"time"

	"paepcke.de/gpstime/gpsfeed"
)

//
// SIMPLE API
//

// CheckSetClock anligns continusly the local system to to gpsDevice time
func CheckSetClock(device string) { checkSetClock(&gpsfeed.GpsDevice{FileIO: device}) }

// CheckSetClockOnce sets the clock to gpsDevice time and exists
func CheckSetClockOnce(device string) { checkSetClockOnce(&gpsfeed.GpsDevice{FileIO: device}) }

// GetTime get a single timestamp
func GetTime(device string) time.Time { return parseTime(&gpsfeed.GpsDevice{FileIO: device}) }

// GetLocation gets single location stamp (lat, long, elevation)
func GetLocation(device string) Coord {
	return parseLocation(&gpsfeed.GpsDevice{FileIO: device})
}

// WriteLocationFile reads gps nmea frames from a gps device and writes the current location as unix shell env include
func WriteLocationFile(device string) { writeLocationFile(&gpsfeed.GpsDevice{FileIO: device}) }

// CheckSetClockLocation updates the local system time to gpsDevice time and location env file ( when needed ) in a loop
func CheckSetClockLocation(device string) { checkSetClockLocation(&gpsfeed.GpsDevice{FileIO: device}) }

//
// GENERIC BACKEND
//

// Coord is an gps coordinates point
type Coord struct {
	Valid                bool
	Lat, Long, Elevation float64
}

// CheckSetClockD ...
func CheckSetClockD(dev *gpsfeed.GpsDevice) { checkSetClock(dev) }

// CheckSetClockOnceD ...
func CheckSetClockOnceD(dev *gpsfeed.GpsDevice) { checkSetClockOnce(dev) }

// GetTimeD ...
func GetTimeD(dev *gpsfeed.GpsDevice) time.Time { return parseTime(dev) }

// GetLocationD ...
func GetLocationD(dev *gpsfeed.GpsDevice) Coord { return parseLocation(dev) }

// WriteLocationFileD ...
func WriteLocationFileD(dev *gpsfeed.GpsDevice) { writeLocationFile(dev) }

// WriteLocationFileOnceD ...
func WriteLocationFileOnceD(dev *gpsfeed.GpsDevice) { writeLocationFileOnce(dev) }

// CheckSetClockLocationD ...
func CheckSetClockLocationD(dev *gpsfeed.GpsDevice) { checkSetClockLocation(dev) }

//
// LITTLE HELPER
//

// Out ...
func Out(msg string) { out(msg) }
