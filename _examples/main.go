package main

//go:generate buildcond --pkg main --output . mydebug

import (
	"fmt"

	"bitbucket.org/shu/buildcond/debug"
	"bitbucket.org/shu/buildcond/test"
)

func main() {
	debug.Do(func() {
		fmt.Println("debug")
	})
	Do(func() {
		fmt.Println("MY debug")
	})
	test.Do(func() {
		fmt.Println("testing!!??")
	})
}
