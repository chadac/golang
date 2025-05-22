// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

type s[T any] struct {
	a T
}

func (x s[T]) f() T {
	return x.a
}
func main() {
	x := s[int]{a: 7}
	f := x.f
	if golangt, want := f(), 7; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
}
