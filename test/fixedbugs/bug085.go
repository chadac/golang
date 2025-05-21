// errorcheck

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package P

var x int

func foo() {
	print(P.x);  // ERROR "undefined"
}

/*
uetli:~/Source/golang1/test/bugs gri$ 6g bug085.golang
bug085.golang:6: P: undefined
Bus error
*/

/* expected scope hierarchy (outermost to innermost)

universe scope (contains predeclared identifiers int, float32, int32, len, etc.)
"solar" scope (just holds the package name P so it can be found but doesn't conflict)
global scope (the package global scope)
local scopes (function scopes)
*/
