// compile

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Gccgolang used to incorrectly give an error when compiling this.

package p

func F() (i int) {
	for first := true; first; first = false {
		i++
	}
	return
}
