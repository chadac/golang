// run

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Make sure address calculations don't float up before
// the corresponding nil check.

package main

type T struct {
	a, b int
}

//golang:noinline
func f(x *T, p *bool, n int) {
	*p = n != 0
	useStack(1000)
	g(&x.b)
}

//golang:noinline
func g(p *int) {
}

func useStack(n int) {
	if n == 0 {
		return
	}
	useStack(n - 1)
}

func main() {
	mustPanic(func() {
		var b bool
		f(nil, &b, 3)
	})
}

func mustPanic(f func()) {
	defer func() {
		if recover() == nil {
			panic("expected panic, golangt nil")
		}
	}()
	f()
}
