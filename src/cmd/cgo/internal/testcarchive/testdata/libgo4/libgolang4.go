// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
#include <signal.h>
#include <pthread.h>

// Raise SIGIO.
static void CRaiseSIGIO(pthread_t* p) {
	pthread_kill(*p, SIGIO);
}
*/
import "C"

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

var sigioCount int32

// Catch SIGIO.
//
//export GolangCatchSIGIO
func GolangCatchSIGIO() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGIO)
	golang func() {
		for range c {
			atomic.AddInt32(&sigioCount, 1)
		}
	}()
}

// Raise SIGIO.
//
//export GolangRaiseSIGIO
func GolangRaiseSIGIO(p *C.pthread_t) {
	C.CRaiseSIGIO(p)
}

// Return the number of SIGIO signals seen.
//
//export SIGIOCount
func SIGIOCount() C.int {
	return C.int(atomic.LoadInt32(&sigioCount))
}

func main() {
}
