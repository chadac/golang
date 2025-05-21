// compile

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Crashed gccgolang.

package p

var F func([0]int) int
var G func() [0]int

var V = make([]int, F(G()))
