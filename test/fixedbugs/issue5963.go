// run

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Used to die in runtime due to init golangroutine exiting while
// locked to main thread.

package main

import (
	"os"
	"runtime"
)

func init() {
	c := make(chan int, 1)
	defer func() {
		c <- 0
	}()
	golang func() {
		os.Exit(<-c)
	}()
	runtime.Goexit()
}

func main() {
}

/* Before fix:

invalid m->locked = 2
fatal error: internal lockOSThread error

golangroutine 2 [runnable]:
runtime.MHeap_Scavenger()
	/Users/rsc/g/golang/src/pkg/runtime/mheap.c:438
runtime.golangexit()
	/Users/rsc/g/golang/src/pkg/runtime/proc.c:1313
created by runtime.main
	/Users/rsc/g/golang/src/pkg/runtime/proc.c:165

golangroutine 3 [runnable]:
main.func·002()
	/Users/rsc/g/golang/test/fixedbugs/issue5963.golang:22
created by main.init·1
	/Users/rsc/g/golang/test/fixedbugs/issue5963.golang:24 +0xb9
exit status 2
*/
