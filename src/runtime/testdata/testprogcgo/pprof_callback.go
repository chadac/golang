// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !plan9 && !windows

package main

// Make many C-to-Go callback while collecting a CPU profile.
//
// This is a regression test for issue 50936.

/*
#include <unistd.h>

void golangCallbackPprof();

static void callGo() {
	// Spent >20us in C so this thread is eligible for sysmon to retake its
	// P.
	usleep(50);
	golangCallbackPprof();
}
*/
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func init() {
	register("CgolangPprofCallback", CgolangPprofCallback)
}

//export golangCallbackPprof
func golangCallbackPprof() {
	// No-op. We want to stress the cgolangcall and cgolangcallback internals,
	// landing as many pprof signals there as possible.
}

func CgolangPprofCallback() {
	// Issue 50936 was a crash in the SIGPROF handler when the signal
	// arrived during the exitsyscall following a cgolangcall(back) in dropg or
	// execute, when updating mp.curg.
	//
	// These are reachable only when exitsyscall finds no P available. Thus
	// we make C calls from significantly more Gs than there are available
	// Ps. Lots of runnable work combined with >20us spent in callGo makes
	// it possible for sysmon to retake Ps, forcing C calls to golang down the
	// desired exitsyscall path.
	//
	// High GOMAXPROCS is used to increase opportunities for failure on
	// high CPU machines.
	const (
		P = 16
		G = 64
	)
	runtime.GOMAXPROCS(P)

	f, err := os.CreateTemp("", "prof")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	for i := 0; i < G; i++ {
		golang func() {
			for {
				C.callGo()
			}
		}()
	}

	time.Sleep(time.Second)

	pprof.StopCPUProfile()

	fmt.Println("OK")
}
