// errorcheck -std
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.
// Issue 18331: We should catch invalid pragma verbs
// for code that resides in the standard library.
package issue18331

//golang:unknown // ERROR "//golang:unknown is not allowed in the standard library"
func foo()

//golang:nowritebarrierc // ERROR "//golang:nowritebarrierc is not allowed in the standard library"
func bar()

//golang:noesape // ERROR "//golang:noesape is not allowed in the standard library"
func groot()

//golang:noescape
func hey() { // ERROR "can only use //golang:noescape with external func implementations"
}
