package cond

import "flag"

var testing bool

func init() {
	testing = (flag.Lookup("test.v") != nil)
}

// IfTest executes function f if in `go test`.
func IfTest(f func()) {
	if testing {
		f()
	}
}

// UnlessTest executes function f if not in `go test`.
func UnlessTest(f func()) {
	if !testing {
		f()
	}
}
