// build

//golang:build cgolang

// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
#cgolang CFLAGS: -Werror -Wimplicit-function-declaration

#include <stdio.h>

static void CFn(_GoString_ golangstr) {
	printf("%.*s\n", (int)(_GoStringLen(golangstr)), _GoStringPtr(golangstr));
}
*/
import "C"

func main() {
	C.CFn("hello, world")
}

// The bug only occurs if there is an exported function.
//export Fn
func Fn() {
}
