// errorcheck

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that misplaced directives are diagnosed.

//golang:build wasm

package p

//golang:wasmexport F
func F() {} // OK

type S int32

//golang:wasmexport M
func (S) M() {} // ERROR "cannot use //golang:wasmexport on a method"
