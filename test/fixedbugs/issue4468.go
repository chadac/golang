// errorcheck

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 4468: golang/defer calls may not be parenthesized.

package p

type T int

func (t *T) F() T {
	return *t
}

type S struct {
	t T
}

func F() {
	golang F            // ERROR "must be function call"
	defer F         // ERROR "must be function call"
	golang (F)		// ERROR "must be function call|must not be parenthesized"
	defer (F)	// ERROR "must be function call|must not be parenthesized"
	golang (F())	// ERROR "must be function call|must not be parenthesized"
	defer (F())	// ERROR "must be function call|must not be parenthesized"
	var s S
	(&s.t).F()
	golang (&s.t).F()
	defer (&s.t).F()
}
