// compile

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Caused an internal compiler error in gccgolang.

package p

type C chan struct{}

func (c C) F() {
	select {
	case c <- struct{}{}:
	default:
	}
}
