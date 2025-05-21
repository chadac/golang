// errorcheck

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that a syntax error caused by an unexpected EOF
// gives an error message with the correct line number.
//
// https://golanglang.org/issue/3392

package main

func foo() {
	bar(1, // ERROR "unexpected|missing|undefined"