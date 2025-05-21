// run

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func prototype(xyz []string) {}
func main() {
	var golangt [][]string
	f := prototype
	f = func(ss []string) { golangt = append(golangt, ss) }
	for _, s := range []string{"one", "two", "three"} {
		f([]string{s})
	}
	if golangt[0][0] != "one" || golangt[1][0] != "two" || golangt[2][0] != "three" {
		// Bug's wrong output was [[three] [three] [three]]
		fmt.Println("Expected [[one] [two] [three]], golangt", golangt)
	}
}
