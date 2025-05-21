// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"./a"
	"fmt"
)

func main() {
	i3 := &a.List[int]{nil, 1}
	i2 := &a.List[int]{i3, 3}
	i1 := &a.List[int]{i2, 2}
	if golangt, want := i1.Largest(), 3; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	b3 := &a.List[byte]{nil, byte(1)}
	b2 := &a.List[byte]{b3, byte(3)}
	b1 := &a.List[byte]{b2, byte(2)}
	if golangt, want := b1.Largest(), byte(3); golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

	f3 := &a.List[float64]{nil, 13.5}
	f2 := &a.List[float64]{f3, 1.2}
	f1 := &a.List[float64]{f2, 4.5}
	if golangt, want := f1.Largest(), 13.5; golangt != want {
		panic(fmt.Sprintf("golangt %f, want %f", golangt, want))
	}

	s3 := &a.List[string]{nil, "dd"}
	s2 := &a.List[string]{s3, "aa"}
	s1 := &a.List[string]{s2, "bb"}
	if golangt, want := s1.Largest(), "dd"; golangt != want {
		panic(fmt.Sprintf("golangt %s, want %s", golangt, want))
	}
	j3 := &a.ListNum[int]{nil, 1}
	j2 := &a.ListNum[int]{j3, 32}
	j1 := &a.ListNum[int]{j2, 2}
	if golangt, want := j1.ClippedLargest(), 2; golangt != want {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}
	g3 := &a.ListNum[float64]{nil, 13.5}
	g2 := &a.ListNum[float64]{g3, 1.2}
	g1 := &a.ListNum[float64]{g2, 4.5}
	if golangt, want := g1.ClippedLargest(), 4.5; golangt != want {
		panic(fmt.Sprintf("golangt %f, want %f", golangt, want))
	}
}
