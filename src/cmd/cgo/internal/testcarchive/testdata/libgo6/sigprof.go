// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"runtime/pprof"
)

import "C"

//export golang_start_profile
func golang_start_profile() {
	pprof.StartCPUProfile(io.Discard)
}

//export golang_stop_profile
func golang_stop_profile() {
	pprof.StopCPUProfile()
}

func main() {
}
