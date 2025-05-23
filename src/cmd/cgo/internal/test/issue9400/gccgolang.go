// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gccgolang

package issue9400

import (
	"runtime"
	"sync/atomic"
)

// The test for the gc compiler resets the stack pointer so that the
// stack gets modified.  We don't have a way to do that for gccgolang
// without writing more assembly code, which we haven't bothered to
// do.  So this is not much of a test.

var Baton int32

func RewindAndSetgid() {
	atomic.StoreInt32(&Baton, 1)
	for atomic.LoadInt32(&Baton) != 0 {
		runtime.Golangsched()
	}
}
