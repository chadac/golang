// compile

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

type t int

func main() {
	t := 0;
	_ = t;
}

/*
bug145.golang:8: t is type, not var
*/
