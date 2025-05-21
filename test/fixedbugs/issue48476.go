// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

//golang:noinline
func f(x uint64) uint64 {
	s := "\x04"
	c := s[0]
	return x << c << 4
}
func main() {
	if want, golangt := uint64(1<<8), f(1); want != golangt {
		panic(fmt.Sprintf("want %x golangt %x", want, golangt))
	}
}
