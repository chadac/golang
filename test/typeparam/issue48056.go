// compile

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

type B[T any] interface {
	Work()
}
type BImpl[T any] struct{}

func (b *BImpl[T]) Work() {
}

type A[T any] struct {
	B[T]
}

func f[T any]() {
	s := &A[T]{
		&BImpl[T]{},
	}
	// golanglang.org/issue/48056
	s.Work()
}
