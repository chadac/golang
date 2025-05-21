// compile

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test case for issue 475. This file should compile.

package main

import . "unsafe"

func main() {
	var x int
	println(Sizeof(x))
}

/*
bug239.golang:11: imported and not used: unsafe
bug239.golang:15: undefined: Sizeof
*/
