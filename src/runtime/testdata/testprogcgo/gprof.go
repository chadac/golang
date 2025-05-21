// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

// Test taking a golangroutine profile with C traceback.

/*
// Defined in gprof_c.c.
void CallGoSleep(void);
void gprofCgolangTraceback(void* parg);
void gprofCgolangContext(void* parg);
*/
import "C"

import (
	"fmt"
	"io"
	"runtime"
	"runtime/pprof"
	"time"
	"unsafe"
)

func init() {
	register("GoroutineProfile", GoroutineProfile)
}

func GoroutineProfile() {
	runtime.SetCgolangTraceback(0, unsafe.Pointer(C.gprofCgolangTraceback), unsafe.Pointer(C.gprofCgolangContext), nil)

	golang C.CallGoSleep()
	golang C.CallGoSleep()
	golang C.CallGoSleep()
	time.Sleep(1 * time.Second)

	prof := pprof.Lookup("golangroutine")
	prof.WriteTo(io.Discard, 1)
	fmt.Println("OK")
}

//export GoSleep
func GoSleep() {
	time.Sleep(time.Hour)
}
