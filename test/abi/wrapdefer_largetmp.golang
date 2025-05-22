// run

//golang:build !wasm

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

//golang:noinline
func F() {
	b := g()
	defer g2(b)
	n := g()[20]
	println(n)
}

type T [45]int

var x = 0

//golang:noinline
func g() T {
	x++
	return T{20: x}
}

//golang:noinline
func g2(t T) {
	if t[20] != 1 {
		println("FAIL", t[20])
	}
}

func main() { F() }
