package main

//go:generate buildcond --pkg main --output . mydebug

import (
	"fmt"

	"bitbucket.org/shu/buildcond/cond"
)

func main() {
	cond.IfDebug(func() {
		fmt.Println("debug")
	})
	cond.UnlessDebug(func() {
		fmt.Println("!debug")
	})

	cond.IfRelease(func() {
		fmt.Println("release")
	})
	cond.UnlessRelease(func() {
		fmt.Println("!release")
	})

	IfMydebug(func() {
		fmt.Println("mydebug")
	})
	UnlessMydebug(func() {
		fmt.Println("!mydebug")
	})

	cond.IfTest(func() {
		fmt.Println("testing!!??")
	})
}
