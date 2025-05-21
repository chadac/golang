// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	var s string;
	s = "0000000000000000000000000000000000000000000000000000000000"[0:7];
	_ = s;
}

/*
uetli:~/Source/golang1/test/bugs gri$ 6g bug061.golang
Bus error
*/
