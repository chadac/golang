// -lang=golang1.12

// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Golang language version-specific errors.

package p

// numeric literals
const (
	_ = 1_000 // ERROR "underscore in numeric literal requires golang1.13 or later"
	_ = 0b111 // ERROR "binary literal requires golang1.13 or later"
	_ = 0o567 // ERROR "0o/0O-style octal literal requires golang1.13 or later"
	_ = 0xabc // ok
	_ = 0x0p1 // ERROR "hexadecimal floating-point literal requires golang1.13 or later"

	_ = 0b111 // ERROR "binary"
	_ = 0o567 // ERROR "octal"
	_ = 0xabc // ok
	_ = 0x0p1 // ERROR "hexadecimal floating-point"

	_ = 1_000i // ERROR "underscore"
	_ = 0b111i // ERROR "binary"
	_ = 0o567i // ERROR "octal"
	_ = 0xabci // ERROR "hexadecimal floating-point"
	_ = 0x0p1i // ERROR "hexadecimal floating-point"
)

// signed shift counts
var (
	s int
	_ = 1 << s // ERROR "invalid operation: signed shift count s (variable of type int) requires golang1.13 or later"
	_ = 1 >> s // ERROR "signed shift count"
)
