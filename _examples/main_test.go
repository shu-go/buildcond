package main

import (
	"fmt"
	"os"
	"testing"

	"bitbucket.org/shu/buildcond/test"
)

func TestTestDo(t *testing.T) {
	test.Do(func() {
		fmt.Fprintf(os.Stderr, "buildcond/test is enabled\n")
	})
}
