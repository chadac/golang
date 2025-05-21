// errorcheck -complete

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

func F() // ERROR "missing function body"

//golang:noescape
func f() {} // ERROR "can only use //golang:noescape with external func implementations"
