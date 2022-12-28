# OVERVIEW

- set local system time from usb gps dongle (nmea mode)
- set local site coord information from usb gps dongle 
- focus on embedded systems, enviromental friendly systems
- focus on low energy use and secure and relyable operation
- not the right choice for sub second critical applications (!)
- 100% pure go, minimal(internal-only) imports
- use as app or api (see api.go)

# INSTALL

```
go install paepcke.de/gpstime/cmd/gpstime@latest
```

# SHOWTIME

## set time 
``` Shell 
gpstime
Apr 25 06:50:25 rpi2b32-pnoc gpstime[49152]:  ### INFORMAL ### LARGE [initial] TIMEJUMP ### TRYING FAST FORWARD 
Apr 25 06:53:02 rpi2b32-pnoc gpstime[49152]: srv  time: 2022-04-25 06:53:02 +0000 UTC 
Apr 25 06:53:02 rpi2b32-pnoc gpstime[49152]: loc  time: 2022-04-25 06:50:25.734559739 +0000 UTC m=+1.230512813 
Apr 25 06:53:02 rpi2b32-pnoc gpstime[49152]: diff time: 2m36.265440261s 
```

## set localtion
``` Shell 
gpstime location /dev/gps0
cat /var/gps/.location]
#!/bin/sh
export GPS_MODE="gpstime"
export GPS_LAT="53.56409"
export GPS_LONG="9.95747"
export GPS_ELEVATION="0"
export AIRLOCTAG="128@X3UR-KY-LURRJI-UJ-DXHY"
export GPS_SUN_RISE="03:02:02"
export GPS_SUN_SET="19:33:42"
export GPS_SUN_NOON="11:17:24"
export GPS_SUN_DAYLIGHT="16h31m40s"
```

# CONTRIBUTION

Yes, Please! PRs Welcome! 
