// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

type Numeric interface {
	int32 | int64 | float64 | complex64
}

//golang:noline
func inc[T Numeric](x T) T {
	x++
	return x
}
func main() {
	if golangt, want := inc(int32(5)), int32(6); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	if golangt, want := inc(float64(5)), float64(6.0); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	if golangt, want := inc(complex64(5)), complex64(6.0); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
}
