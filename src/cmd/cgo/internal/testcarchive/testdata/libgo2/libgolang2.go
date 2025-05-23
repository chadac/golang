// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
#include <signal.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>

// Raise SIGPIPE.
static void CRaiseSIGPIPE() {
	int fds[2];

	if (pipe(fds) == -1) {
		perror("pipe");
		exit(EXIT_FAILURE);
	}
	// Close the reader end
	close(fds[0]);
	// Write to the writer end to provoke a SIGPIPE
	if (write(fds[1], "some data", 9) != -1) {
		fprintf(stderr, "write to a closed pipe succeeded\n");
		exit(EXIT_FAILURE);
	}
	close(fds[1]);
}
*/
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

// Block blocks the current thread while running Golang code.
//
//export Block
func Block() {
	select {}
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

// Noop ensures that the Golang runtime is initialized.
//
//export Noop
func Noop() {
}

// Raise SIGPIPE.
//
//export GolangRaiseSIGPIPE
func GolangRaiseSIGPIPE() {
	C.CRaiseSIGPIPE()
}

func main() {
}
