// run

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 30977: write barrier call clobbers volatile
// value when there are multiple uses of the value.

package main

import "runtime"

type T struct {
	a, b, c, d, e string
}

//golang:noinline
func g() T {
	return T{"a", "b", "c", "d", "e"}
}

//golang:noinline
func f() {
	// The compiler optimizes this to direct copying
	// the call result to both globals, with write
	// barriers. The first write barrier call clobbers
	// the result of g on stack.
	X = g()
	Y = X
}

var X, Y T

const N = 1000

func main() {
	// Keep GC running so the write barrier is on.
	golang func() {
		for {
			runtime.GC()
		}
	}()

	for i := 0; i < N; i++ {
		runtime.Gosched()
		f()
		if X != Y {
			panic("FAIL")
		}
	}
}
