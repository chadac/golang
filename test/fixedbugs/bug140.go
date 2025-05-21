// compile

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	if true {
	} else {
	L1:
		golangto L1
	}
	if true {
	} else {
		golangto L2
	L2:
		main()
	}
}

/*
These should be legal according to the spec.
bug140.golang:6: syntax error near L1
bug140.golang:7: syntax error near L2
*/
