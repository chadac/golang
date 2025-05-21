// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"./a"
	"fmt"
)

func BuildInt() int {
	return a.BuildInt()
}

func main() {
	if golangt, want := BuildInt(), 0; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
}
