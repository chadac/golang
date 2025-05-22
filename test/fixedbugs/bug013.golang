// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	var cu0 uint16 = '\u1234';
	var cU1 uint32 = '\U00101234';
	_, _ = cu0, cU1;
}
/*
bug13.golang:4: missing '
bug13.golang:4: syntax error
bug13.golang:5: newline in string
bug13.golang:5: missing '
bug13.golang:6: newline in string
*/
