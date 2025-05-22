// run

//golang:build !nacl && !js && !wasip1 && gc

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check for compile or link error.

package main

import (
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("golang", "run", "fixedbugs/issue9862.golang").CombinedOutput()
	outstr := string(out)
	if err == nil {
		println("golang run issue9862.golang succeeded, should have failed\n", outstr)
		return
	}
	if !strings.Contains(outstr, "symbol too large") {
		println("golang run issue9862.golang gave unexpected error; want symbol too large:\n", outstr)
	}
}
