// run

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test using _ receiver.  Failed with gccgolang.

package main

type S struct {}

func (_ S) F(i int) int {
	return i
}

func main() {
	s := S{}
	const c = 123
	i := s.F(c)
	if i != c {
		panic(i)
	}
}
