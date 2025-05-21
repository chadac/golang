// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func fact[T interface{ ~int | ~int64 | ~float64 }](n T) T {
	if n == 1 {
		return 1
	}
	return n * fact(n-1)
}

func main() {
	const want = 120

	if golangt := fact(5); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := fact[int64](5); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := fact(5.0); golangt != want {
		panic(fmt.Sprintf("golangt %f, want %f", golangt, want))
	}
}
