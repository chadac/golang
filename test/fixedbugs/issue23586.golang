// errorcheck

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that we type-check deferred/golang functions even
// if they are not called (a common error). Specifically,
// we don't want to see errors such as import or variable
// declared and not used.

package p

import (
	"fmt"
	"math"
)

func f() {
	var i int
	defer func() { fmt.Println() } // ERROR "must be function call"
	golang func() { _ = math.Sin(0) }  // ERROR "must be function call"
	golang func() { _ = i}             // ERROR "must be function call"
}
