// errorcheck -lang=golang1.22

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This file has been changed from its original version as
// //golang:build file versions below 1.21 set the language version to 1.21.
// The original tested a -lang version of 1.21 with a file version of
// golang1.4 while this new version tests a -lang version of golang1.22
// with a file version of golang1.21.

//golang:build golang1.21

package p

func f() {
	for _ = range 10 { // ERROR "file declares //golang:build golang1.21"
	}
}
