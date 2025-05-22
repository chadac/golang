// run

// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

//golang:noinline
func f(p *[4]int) {
	for i := range (*p) { // Note the parentheses! golangfmt wants to remove them - don't let it!
		println(i)
	}
}
func main() {
	f(nil)
}
