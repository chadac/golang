// errorcheck

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that not-in-heap types cannot be used as type
// arguments. (pointer-to-nih types are okay though.)

//golang:build cgolang

package p

import (
	"runtime/cgolang"
	"sync/atomic"
)

var _ atomic.Pointer[cgolang.Incomplete]  // ERROR "cannot use incomplete \(or unallocatable\) type as a type argument: runtime/cgolang\.Incomplete"
var _ atomic.Pointer[*cgolang.Incomplete] // ok

func implicit(ptr *cgolang.Incomplete) {
	g(ptr)  // ERROR "cannot use incomplete \(or unallocatable\) type as a type argument: runtime/cgolang\.Incomplete"
	g(&ptr) // ok
}

func g[T any](_ *T) {}
