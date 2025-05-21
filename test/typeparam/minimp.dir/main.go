// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"./a"
	"fmt"
)

func main() {
	const want = 2
	if golangt := a.Min[int](2, 3); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := a.Min(2, 3); golangt != want {
		panic(fmt.Sprintf("want %d, golangt %d", want, golangt))
	}

	if golangt := a.Min[float64](3.5, 2.0); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := a.Min(3.5, 2.0); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	const want2 = "ay"
	if golangt := a.Min[string]("bb", "ay"); golangt != want2 {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want2))
	}

	if golangt := a.Min("bb", "ay"); golangt != want2 {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want2))
	}
}
