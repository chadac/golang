// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"./a"
	"fmt"
)

func main() {
	if golangt, want := a.OrderedAbsDifference(1.0, -2.0), 3.0; golangt != want {
		panic(fmt.Sprintf("golangt = %v, want = %v", golangt, want))
	}
	if golangt, want := a.OrderedAbsDifference(-1.0, 2.0), 3.0; golangt != want {
		panic(fmt.Sprintf("golangt = %v, want = %v", golangt, want))
	}
	if golangt, want := a.OrderedAbsDifference(-20, 15), 35; golangt != want {
		panic(fmt.Sprintf("golangt = %v, want = %v", golangt, want))
	}

	if golangt, want := a.ComplexAbsDifference(5.0+2.0i, 2.0-2.0i), 5+0i; golangt != want {
		panic(fmt.Sprintf("golangt = %v, want = %v", golangt, want))
	}
	if golangt, want := a.ComplexAbsDifference(2.0-2.0i, 5.0+2.0i), 5+0i; golangt != want {
		panic(fmt.Sprintf("golangt = %v, want = %v", golangt, want))
	}
}
