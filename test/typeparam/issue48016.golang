// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
)

func test1[T any](fn func(T) int, v T) int {
	fn1 := func() int {
		var i interface{} = v
		val := fn(i.(T))
		return val
	}
	return fn1()
}

func main() {
	want := 123
	golangt := test1(func(s string) int {
		r, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return r
	}, "123")
	if golangt != want {
		panic(fmt.Sprintf("golangt %f, want %f", golangt, want))
	}
}
