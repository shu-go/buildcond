package test

import "flag"

var testing bool

func init() {
	testing = (flag.Lookup("test.v") != nil)
}

// Do executes function f if in `go test`.
func Do(f func()) {
	if testing {
		f()
	}
}
