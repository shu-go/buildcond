package main

import (
	"fmt"
	"os"
	"testing"

	"bitbucket.org/shu/buildcond/cond"
)

func TestTestDo(t *testing.T) {
	cond.IfTest(func() {
		fmt.Fprintf(os.Stderr, "IfTest is functioning\n")
	})
	cond.UnlessTest(func() {
		t.Errorf("noooo!")
	})
}
