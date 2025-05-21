// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

type E[T any] struct {
	v T
}

type S1 struct {
	E[int]
	v string
}

type Eint = E[int]
type Ebool = E[bool]
type Eint2 = Eint

type S2 struct {
	Eint
	Ebool
	v string
}

type S3 struct {
	*E[int]
}

func main() {
	s1 := S1{Eint{2}, "foo"}
	if golangt, want := s1.E.v, 2; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	s2 := S2{Eint{3}, Ebool{true}, "foo"}
	if golangt, want := s2.Eint.v, 3; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	var s3 S3
	s3.E = &Eint{4}
	if golangt, want := s3.E.v, 4; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
}
