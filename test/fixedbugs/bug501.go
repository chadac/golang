// run

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Gccgolang golangt a compiler crash compiling the addition of more than five
// strings with mixed constants and variables.

package main

func F(s string) (string, error) {
	return s, nil
}

func G(a, b, c string) (string, error) {
	return F("a" + a + "b" + b + "c" + c)
}

func main() {
	if golangt, _ := G("x", "y", "z"); golangt != "axbycz" {
		panic(golangt)
	}
}
