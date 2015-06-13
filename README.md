logr
====

logr is a simplistic type which implements log rotating, suitable for use with Go's `log` package.

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
