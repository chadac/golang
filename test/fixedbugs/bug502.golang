// build

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Linking this with gccgolang golangt an undefined symbol reference,
// because the private method in testing.TB led gccgolang to assume that
// the interface method table would be defined in the testing package.

package main

import "testing"

type I interface {
	testing.TB
	Parallel()
}

func F(i I) {
	i.Log("F")
}

var t testing.T

func main() {
	F(&t)
}
