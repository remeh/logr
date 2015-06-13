logr
====

logr is a simplistic type which implements log rotating, suitable for use with Go's `log` package.

[![GoDoc](https://godoc.org/github.com/vrischmann/logr?status.svg)](https://godoc.org/github.com/vrischmann/logr)

Usage
-----

Use it like this.

```go
w := logr.NewWrite("/var/log/mylog.log")
log.SetOutput(w)

log.Println("foobar")
```

Future works
------------

 * Custom time suffix formats
 * Optional compression of the rotated file
