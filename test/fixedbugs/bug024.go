// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	var i int;
	i = '\'';
	i = '\\';
	var s string;
	s = "\"";
	_, _ = i, s;
}
/*
bug.golang:5: unknown escape sequence: '
bug.golang:6: unknown escape sequence: \
bug.golang:8: unknown escape sequence: "
*/
