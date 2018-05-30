// +build release

package cond

// IfRelease executes function f if the build tag 'release' is enabled.
func IfRelease(f func()) {
	f()
}

// UnlessRelease executes function f if the build tag 'release' is disabled.
func UnlessRelease(f func()) {
	// f()
}
