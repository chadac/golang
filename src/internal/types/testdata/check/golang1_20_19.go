// -lang=golang1.20

// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Golang language version-specific errors.

//golang:build golang1.19

package p

type Slice []byte
type Array [8]byte

var s Slice
var p = (Array)(s /* ok because file versions below golang1.21 set the language version to golang1.21 */)
