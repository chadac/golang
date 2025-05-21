// compile

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

const s string = "foo";

func main() {
	i := len(s);  // should be legal to take len() of a constant
	_ = i;
}

/*
uetli:~/Source/golang1/test/bugs gri$ 6g bug087.golang
bug087.golang:6: illegal combination of literals LEN 9
bug087.golang:6: illegal combination of literals LEN 9
*/
