// Package buildcond provides a functionality to detect a build tag.
//
// This package provides following go functions:
//   buildcond/cond/IfDebug, UnlessDebug
//   buildcond/cond/IfRelease, UnlessRelease
//   buildcond/cond/IfTest, UnlessTest
//
// IfTest is for `go build -tags "debug"`, and IfRelease is for `go build -tags "release"`.
//
// IfTest is for `go test`, not for `go build -tags "test"`.
//
// An executable `buildcond` (src is in cmd/buildcond) is a go-generate command.
// It generates {tag}.go and no{tag}.go, those have a function Do(f func()) to execute f if the buildtag is enabled.
// See `buildcond -h` for details and usages.
//
package buildcond

//go:generate buildcond --pkg=cond debug
//go:generate buildcond --pkg=cond release
