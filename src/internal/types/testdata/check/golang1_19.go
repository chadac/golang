// -lang=golang1.19

// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Golang language version-specific errors.

package p

type Slice []byte
type Array [8]byte

var s Slice
var p = (Array)(s /* ERROR "requires golang1.20 or later" */)
