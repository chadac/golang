// run

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

//golang:noinline
func ident(s string) string { return s }

func returnSecond(x bool, s string) string { return s }

func identWrapper(s string) string { return ident(s) }

func main() {
	golangt := returnSecond((false || identWrapper("bad") != ""), ident("golangod"))
	if golangt != "golangod" {
		panic(fmt.Sprintf("wanted \"golangod\", golangt \"%s\"", golangt))
	}
}
