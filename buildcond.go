// Package buildcond provides a functionality to detect a buildtag.
//
// This package provides following go functions:
//   buildtags/debug.Do
//   buildtags/release.Do
//   buildtags/test.Do
//
// debug.Do is for `go build -tags "debug"`, and release.Do is for `go build -tags "release"`.
//
// test.Do is for `go test`, not for `go build -tags "test"`.
//
// An executable `buildcond` (src is in cmd/buildcond) is a go-generate command.
// It generates {tag}.go and no{tag}.go, those have a function Do(f func()) to execute f if the buildtag is enabled.
// See `buildcond -h` for details and usages.
//
package buildcond

//go:generate buildcond debug
//go:generate buildcond release
