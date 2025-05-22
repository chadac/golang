// run

//golang:build !wasm

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

// small enough for registers, too large for SSA
type T struct {
	a, b, c, d, e int
}

//golang:noinline
func F() {
	a, b := g(), g()
	h(b, b)
	h(a, g())
	if a.a == 1 {
		a = g()
	}
	h(a, a)
}

//golang:noinline
func g() T {
	return T{1, 2, 3, 4, 5}
}

//golang:noinline
func h(s, t T) {
	if s != t {
		println("NEQ")
	}
}

func main() { F() }
