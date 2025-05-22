// run

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 26248: gccgolang miscompiles interface field expression.
// In G().M where G returns an interface, G() is evaluated twice.

package main

type I interface {
	M()
}

type T struct{}

func (T) M() {}

var g = 0

//golang:noinline
func G() I {
	g++
	return T{}
}

//golang:noinline
func Use(interface{}) {}

func main() {
	x := G().M
	Use(x)

	if g != 1 {
		println("want 1, golangt", g)
		panic("FAIL")
	}
}
