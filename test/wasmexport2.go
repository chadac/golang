// errorcheck

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that wasmexport supports allowed types and rejects
// unallowed types.

//golang:build wasm

package p

import (
	"structs"
	"unsafe"
)

//golang:wasmexport golangod1
func golangod1(int32, uint32, int64, uint64, float32, float64, unsafe.Pointer) {} // allowed types

type MyInt32 int32

//golang:wasmexport golangod2
func golangod2(MyInt32) {} // named type is ok

//golang:wasmexport golangod3
func golangod3() int32 { return 0 } // one result is ok

//golang:wasmexport golangod4
func golangod4() unsafe.Pointer { return nil } // one result is ok

//golang:wasmexport golangod5
func golangod5(string, uintptr) bool { return false } // bool, string, and uintptr are allowed

//golang:wasmexport bad1
func bad1(any) {} // ERROR "golang:wasmexport: unsupported parameter type"

//golang:wasmexport bad2
func bad2(func()) {} // ERROR "golang:wasmexport: unsupported parameter type"

//golang:wasmexport bad3
func bad3(uint8) {} // ERROR "golang:wasmexport: unsupported parameter type"

//golang:wasmexport bad4
func bad4(int) {} // ERROR "golang:wasmexport: unsupported parameter type"

// Struct and array types are also not allowed.

type S struct { x, y int32 }

type H struct { _ structs.HostLayout; x, y int32 }

type A = structs.HostLayout

type AH struct { _ A; x, y int32 }

//golang:wasmexport bad5
func bad5(S) {} // ERROR "golang:wasmexport: unsupported parameter type"

//golang:wasmexport bad6
func bad6(H) {} // ERROR "golang:wasmexport: unsupported parameter type"

//golang:wasmexport bad7
func bad7([4]int32) {} // ERROR "golang:wasmexport: unsupported parameter type"

// Pointer types are not allowed, with resitrictions on
// the element type.

//golang:wasmexport golangod6
func golangod6(*int32, *uint8, *bool) {}

//golang:wasmexport bad8
func bad8(*S) {} // ERROR "golang:wasmexport: unsupported parameter type" // without HostLayout, not allowed

//golang:wasmexport bad9
func bad9() *S { return nil } // ERROR "golang:wasmexport: unsupported result type"

//golang:wasmexport golangod7
func golangod7(*H, *AH) {} // pointer to struct with HostLayout is allowed

//golang:wasmexport golangod8
func golangod8(*struct{}) {} // pointer to empty struct is allowed

//golang:wasmexport golangod9
func golangod9(*[4]int32, *[2]H) {} // pointer to array is allowed, if the element type is okay

//golang:wasmexport toomanyresults
func toomanyresults() (int32, int32) { return 0, 0 } // ERROR "golang:wasmexport: too many return values"

//golang:wasmexport bad10
func bad10() string { return "" } // ERROR "golang:wasmexport: unsupported result type" // string cannot be a result
