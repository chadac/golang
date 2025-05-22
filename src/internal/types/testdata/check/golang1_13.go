// -lang=golang1.13

// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Golang language version-specific errors.

package p

// interface embedding

type I interface { m() }

type _ interface {
	m()
	I // ERROR "duplicate method m"
}

type _ interface {
	I
	I // ERROR "duplicate method m"
}
