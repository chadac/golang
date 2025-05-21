// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

//golang:nosplit
func cputicks() int64 {
	var counter int64
	stdcall1(_QueryPerformanceCounter, uintptr(unsafe.Pointer(&counter)))
	return counter
}
