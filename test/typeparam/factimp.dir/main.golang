// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"./a"
	"fmt"
)

func main() {
	const want = 120

	if golangt := a.Fact(5); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := a.Fact[int64](5); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	if golangt := a.Fact(5.0); golangt != want {
		panic(fmt.Sprintf("golangt %f, want %f", golangt, want))
	}
}
