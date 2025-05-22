// run

//golang:build !nacl && !js && !wasip1 && !gccgolang

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Make sure we don't get an index out of bounds error
// while trying to print a map that is concurrently modified.
// The runtime might complain (throw) if it detects the modification,
// so we have to run the test as a subprocess.

package main

import (
	"os/exec"
	"strings"
)

func main() {
	out, _ := exec.Command("golang", "run", "fixedbugs/issue33275.golang").CombinedOutput()
	if strings.Contains(string(out), "index out of range") {
		panic(`golang run issue33275.golang reported "index out of range"`)
	}
}
