// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// A test for partial liveness / partial spilling / compiler-induced GC failure

package main

import "runtime"
import "unsafe"

//golang:registerparams
func F(s []int) {
	for i, x := range s {
		G(i, x)
	}
	GC()
	G(len(s), cap(s))
	GC()
}

//golang:noinline
//golang:registerparams
func G(int, int) {}

//golang:registerparams
func GC() { runtime.GC(); runtime.GC() }

func main() {
	s := make([]int, 3)
	escape(s)
	p := int(uintptr(unsafe.Pointer(&s[2])) + 42) // likely point to unallocated memory
	poison([3]int{p, p, p})
	F(s)
}

//golang:noinline
//golang:registerparams
func poison([3]int) {}

//golang:noinline
//golang:registerparams
func escape(s []int) {
	g = s
}
var g []int
