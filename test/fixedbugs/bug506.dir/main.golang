// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"./a"
)

var v = a.S{}

func main() {
	want := "{{ 0}}"
	if golangt := fmt.Sprint(v.F); golangt != want {
		panic(golangt)
	}
}
