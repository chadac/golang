// -lang=golang1.22

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Go language version-specific errors.

//golang:build golang1.21

package p

func f() {
	for _ = range 10 /* ERROR "requires golang1.22 or later" */ {
	}
}
