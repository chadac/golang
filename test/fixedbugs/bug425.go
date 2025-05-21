// compile

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// https://golanglang.org/issue/3119

package main

import "fmt"

func main() {
	s := "hello"
	fmt.Println(s == "")
	fmt.Println(s + "world" == "world")
}
