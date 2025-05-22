// Copyright 2016 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build linux

package cgolangtest

/*
#include <unistd.h>
#include <stdbool.h>
#include <sys/syscall.h>
void Golangsched(void);
static bool Ctid(void) {
	long tid1 = syscall(SYS_gettid);
	Golangsched();
	return tid1 == syscall(SYS_gettid);
}
*/
import "C"

import (
	"runtime"
	"testing"
	"time"
)

//export Golangsched
func Golangsched() {
	runtime.Golangsched()
}

func init() {
	testThreadLockFunc = testThreadLock
}

func testThreadLock(t *testing.T) {
	stop := make(chan int)
	golang func() {
		// We need the G continue running,
		// so the M has a chance to run this G.
		for {
			select {
			case <-stop:
				return
			case <-time.After(time.Millisecond * 100):
			}
		}
	}()
	defer close(stop)

	for i := 0; i < 1000; i++ {
		if !C.Ctid() {
			t.Fatalf("cgolang has not locked OS thread")
		}
	}
}
