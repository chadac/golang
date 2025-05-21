// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	var s2 string = "\a\b\f\n\r\t\v";  // \r is miscompiled
	_ = s2;
}
/*
main.golang.c: In function ‘main_main’:
main.golang.c:20: error: missing terminating " character
main.golang.c:21: error: missing terminating " character
main.golang.c:24: error: ‘def’ undeclared (first use in this function)
main.golang.c:24: error: (Each undeclared identifier is reported only once
main.golang.c:24: error: for each function it appears in.)
main.golang.c:24: error: syntax error before ‘def’
main.golang.c:24: error: missing terminating " character
main.golang.c:25: warning: excess elements in struct initializer
main.golang.c:25: warning: (near initialization for ‘slit’)
main.golang.c:36: error: syntax error at end of input
*/
