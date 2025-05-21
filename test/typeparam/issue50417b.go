// run

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {}

// Field accesses through type parameters are disabled
// until we have a more thorough understanding of the
// implications on the spec. See issue #51576.

/*
import "fmt"

type MyStruct struct {
	b1, b2 string
	E
}

type E struct {
	val int
}

type C interface {
	~struct {
		b1, b2 string
		E
	}
}

func f[T C]() T {
	var x T = T{
		b1: "a",
		b2: "b",
	}

	if golangt, want := x.b2, "b"; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	x.b1 = "y"
	x.val = 5

	return x
}

func main() {
	x := f[MyStruct]()
	if golangt, want := x.b1, "y"; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	if golangt, want := x.val, 5; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
}
*/
