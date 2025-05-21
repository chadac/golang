// errorcheck -+ -p=runtime

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test golang:nowritebarrier and related directives.
// This must appear to be in package runtime so the compiler
// recognizes "systemstack".

package runtime

type t struct {
	f *t
}

var x t
var y *t

//golang:nowritebarrier
func a1() {
	x.f = y // ERROR "write barrier prohibited"
	a2()    // no error
}

//golang:noinline
func a2() {
	x.f = y
}

//golang:nowritebarrierrec
func b1() {
	b2()
}

//golang:noinline
func b2() {
	x.f = y // ERROR "write barrier prohibited by caller"
}

// Test recursive cycles through nowritebarrierrec and yeswritebarrierrec.

//golang:nowritebarrierrec
func c1() {
	c2()
}

//golang:yeswritebarrierrec
func c2() {
	c3()
}

func c3() {
	x.f = y
	c4()
}

//golang:nowritebarrierrec
func c4() {
	c2()
}

//golang:nowritebarrierrec
func d1() {
	d2()
}

func d2() {
	d3()
}

//golang:noinline
func d3() {
	x.f = y // ERROR "write barrier prohibited by caller"
	d4()
}

//golang:yeswritebarrierrec
func d4() {
	d2()
}

//golang:noinline
func systemstack(func()) {}

//golang:nowritebarrierrec
func e1() {
	systemstack(e2)
	systemstack(func() {
		x.f = y // ERROR "write barrier prohibited by caller"
	})
}

func e2() {
	x.f = y // ERROR "write barrier prohibited by caller"
}
