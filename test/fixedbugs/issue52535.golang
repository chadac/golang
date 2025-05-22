// compile

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// golangfrontend crashed converting unnamed bool type to any.

package p

func F() {
	m := make(map[int]int)
	var ok any
	_, ok = m[0]
	_ = ok
}
