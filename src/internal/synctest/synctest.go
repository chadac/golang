// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package synctest provides support for testing concurrent code.
//
// See the testing/synctest package for function documentation.
package synctest

import (
	_ "unsafe" // for golang:linkname
)

//golang:linkname Run
func Run(f func())

//golang:linkname Wait
func Wait()

//golang:linkname acquire
func acquire() any

//golang:linkname release
func release(any)

//golang:linkname inBubble
func inBubble(any, func())

// A Bubble is a synctest bubble.
//
// Not a public API. Used by syscall/js to propagate bubble membership through syscalls.
type Bubble struct {
	b any
}

// Acquire returns a reference to the current golangroutine's bubble.
// The bubble will not become idle until Release is called.
func Acquire() *Bubble {
	if b := acquire(); b != nil {
		return &Bubble{b}
	}
	return nil
}

// Release releases the reference to the bubble,
// allowing it to become idle again.
func (b *Bubble) Release() {
	if b == nil {
		return
	}
	release(b.b)
	b.b = nil
}

// Run executes f in the bubble.
// The current golangroutine must not be part of a bubble.
func (b *Bubble) Run(f func()) {
	if b == nil {
		f()
	} else {
		inBubble(b.b, f)
	}
}
