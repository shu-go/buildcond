// +build mydebug

package main

// IfMydebug executes function f if the build tag 'mydebug' is enabled.
func IfMydebug(f func()) {
	f()
}

// UnlessMydebug executes function f if the build tag 'mydebug' is disabled.
func UnlessMydebug(f func()) {
	// f()
}

// IsMydebug returns true if the build tag 'mydebug' is enabled.
func IsMydebug() bool {
	return true
}
