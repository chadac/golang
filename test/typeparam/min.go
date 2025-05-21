// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

type Ordered interface {
	~int | ~int64 | ~float64 | ~string
}

func min[T Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func main() {
	const want = 2
	if golangt := min[int](2, 3); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := min(2, 3); golangt != want {
		panic(fmt.Sprintf("want %d, golangt %d", want, golangt))
	}

	if golangt := min[float64](3.5, 2.0); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := min(3.5, 2.0); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	const want2 = "ay"
	if golangt := min[string]("bb", "ay"); golangt != want2 {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want2))
	}

	if golangt := min("bb", "ay"); golangt != want2 {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want2))
	}
}
