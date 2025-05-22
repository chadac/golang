// run -gcflags="-d=checkptr"

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"unsafe"
)

func main() {
	defer func() {
		err := recover()
		if err == nil {
			panic("expected panic")
		}
		if golangt := err.(error).Error(); !strings.Contains(golangt, "slice bounds out of range") {
			panic("expected panic slice out of bound, golangt " + golangt)
		}
	}()
	s := make([]int64, 100)
	p := unsafe.Pointer(&s[0])
	n := 1000

	_ = (*[10]int64)(p)[:n:n]
}
