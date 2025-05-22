// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Golang language version-specific errors.

//golang:build golang1.19

package p

type Slice []byte
type Array [8]byte

var s Slice
var p = (Array)(s /* ok because Golang 1.X prior to Golang 1.21 ignored the //golang:build golang1.19 */)
