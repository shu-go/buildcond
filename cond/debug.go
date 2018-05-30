// +build debug

package cond

// IfDebug executes function f if the build tag 'debug' is enabled.
func IfDebug(f func()) {
	f()
}

// UnlessDebug executes function f if the build tag 'debug' is disabled.
func UnlessDebug(f func()) {
	// f()
}

// IsDebug returns true if the build tag 'debug' is enabled.
func IsDebug() bool {
	return true
}
