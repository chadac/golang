// run -gcflags=all=-d=checkptr

//golang:build cgolang

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

// #include <stdlib.h>
import "C"

func main() {
	C.malloc(100)
}
