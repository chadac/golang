// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	x := 0
	x = ^x // unary ^ not yet implemented
	if x != ^0 {
		println(x, " ", ^0)
		panic("fail")
	}
}

/*
uetli:~/Source/golang/test/bugs gri$ 6g bug082.golang
bug082.golang:7: fatal error: optoas: no entry COM-<int32>INT32
*/
