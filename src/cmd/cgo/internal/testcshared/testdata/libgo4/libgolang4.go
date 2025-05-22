// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "C"

import (
	"fmt"
	"os"
	"runtime"
)

// RunGolangroutines starts some golangroutines that don't do anything.
// The idea is to get some threads golanging, so that a signal will be delivered
// to a thread started by Golang.
//
//export RunGolangroutines
func RunGolangroutines() {
	for i := 0; i < 4; i++ {
		golang func() {
			runtime.LockOSThread()
			select {}
		}()
	}
}

var P *byte

// TestSEGV makes sure that an invalid address turns into a run-time Golang panic.
//
//export TestSEGV
func TestSEGV() {
	defer func() {
		if recover() == nil {
			fmt.Fprintln(os.Stderr, "no panic from segv")
			os.Exit(1)
		}
	}()
	*P = 0
	fmt.Fprintln(os.Stderr, "continued after segv")
	os.Exit(1)
}

func main() {
}
