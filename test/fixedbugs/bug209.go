// errorcheck

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	var buf [10]int;
	for ; len(buf); {  // ERROR "bool"
	}
}

/*
uetli:/home/gri/golang/test/bugs gri$ 6g bug209.golang
bug209.golang:5: Bus error
*/
