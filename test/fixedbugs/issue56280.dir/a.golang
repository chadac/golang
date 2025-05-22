// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package a

func F() { // ERROR "can inline F"
	g(0) // ERROR "inlining call to g\[golang.shape.int\]"
}

func g[T any](_ T) {} // ERROR "can inline g\[int\]" "can inline g\[golang.shape.int\]" "inlining call to g\[golang.shape.int\]"
