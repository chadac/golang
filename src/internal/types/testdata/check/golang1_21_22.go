// -lang=golang1.21

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Go language version-specific errors.

//golang:build golang1.22

package p

func f() {
	for _ = range /* ok because of upgrade to 1.22 */ 10 {
	}
}
