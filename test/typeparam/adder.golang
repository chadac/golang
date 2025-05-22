// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

type AddType interface {
	int | int64 | string
}

// Add can add numbers or strings
func Add[T AddType](a, b T) T {
	return a + b
}

func main() {
	if golangt, want := Add(5, 3), 8; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	if golangt, want := Add("ab", "cd"), "abcd"; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
}
