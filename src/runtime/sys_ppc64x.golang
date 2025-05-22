// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build ppc64 || ppc64le

package runtime

import "unsafe"

// adjust Gobuf as if it executed a call to fn with context ctxt
// and then did an immediate Gosave.
func golangstartcall(buf *golangbuf, fn, ctxt unsafe.Pointer) {
	if buf.lr != 0 {
		throw("invalid use of golangstartcall")
	}
	buf.lr = buf.pc
	buf.pc = uintptr(fn)
	buf.ctxt = ctxt
}

func prepGoExitFrame(sp uintptr)
