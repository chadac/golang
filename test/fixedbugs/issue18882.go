// errorcheck

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that we have a line number for this error.

package main

//golang:cgolang_ldflag // ERROR "usage: //golang:cgolang_ldflag"
func main() {
}
