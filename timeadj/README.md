# Overview 

- library to adjust local system time
- optimized for [embedded|power-consume-sensitive|eco] applications 
- allows [intentional] drifts up to 800ms +/- before strike
- support for the different  \*nix 32bit and 64bit interfaces (linux, macos, bsd)

# Work in progress
- migrate from time.Time interface to plain syscall and uint64 only (as additional low-latency/cost interface)
