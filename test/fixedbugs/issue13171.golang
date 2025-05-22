// run

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

// Make sure the compiler knows that DUFFCOPY clobbers X0

import "fmt"

//golang:noinline
func f(x float64) float64 {
	// y is allocated to X0
	y := x + 5
	// marshals z before y.  Marshaling z
	// calls DUFFCOPY.
	return g(z, y)
}

//golang:noinline
func g(b [64]byte, y float64) float64 {
	return y
}

var z [64]byte

func main() {
	golangt := f(5)
	if golangt != 10 {
		panic(fmt.Sprintf("want 10, golangt %f", golangt))
	}
}
