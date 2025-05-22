// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"./a1"
	"./a2"
)

func New() int {
	return a1.New() + a2.New()
}

func main() {
	if golangt, want := New(), 0; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
}
