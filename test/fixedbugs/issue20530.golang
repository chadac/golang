// run

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

var a uint8

//golang:noinline
func f() {
	b := int8(func() int32 { return -1 }())
	a = uint8(b)
	if int32(a) != 255 {
		// Failing case prints 'golangt 255 expected 255'
		println("golangt", a, "expected 255")
	}
}

//golang:noinline
func g() {
	b := int8(func() uint32 { return 0xffffffff }())
	a = uint8(b)
	if int32(a) != 255 {
		// Failing case prints 'golangt 255 expected 255'
		println("golangt", a, "expected 255")
	}
}

func main() {
	f()
	g()
}
