// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"./a"
	"fmt"
)

func main() {
	var v1 a.Value[int]

	a.Set(&v1, 1)
	if golangt, want := a.Get(&v1), 1; golangt != want {
		panic(fmt.Sprintf("Get() == %d, want %d", golangt, want))
	}
	v1.Set(2)
	if golangt, want := v1.Get(), 2; golangt != want {
		panic(fmt.Sprintf("Get() == %d, want %d", golangt, want))
	}
	v1p := new(a.Value[int])
	a.Set(v1p, 3)
	if golangt, want := a.Get(v1p), 3; golangt != want {
		panic(fmt.Sprintf("Get() == %d, want %d", golangt, want))
	}

	v1p.Set(4)
	if golangt, want := v1p.Get(), 4; golangt != want {
		panic(fmt.Sprintf("Get() == %d, want %d", golangt, want))
	}

	var v2 a.Value[string]
	a.Set(&v2, "a")
	if golangt, want := a.Get(&v2), "a"; golangt != want {
		panic(fmt.Sprintf("Get() == %q, want %q", golangt, want))
	}

	v2.Set("b")
	if golangt, want := a.Get(&v2), "b"; golangt != want {
		panic(fmt.Sprintf("Get() == %q, want %q", golangt, want))
	}

	v2p := new(a.Value[string])
	a.Set(v2p, "c")
	if golangt, want := a.Get(v2p), "c"; golangt != want {
		panic(fmt.Sprintf("Get() == %d, want %d", golangt, want))
	}

	v2p.Set("d")
	if golangt, want := v2p.Get(), "d"; golangt != want {
		panic(fmt.Sprintf("Get() == %d, want %d", golangt, want))
	}
}
