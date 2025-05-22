// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package one

// Issue 2877
type T struct {
	f func(t *T, arg int)
	g func(t T, arg int)
}

func (t *T) foo(arg int) {}
func (t T) golango(arg int) {}

func (t *T) F() { t.f = (*T).foo }
func (t *T) G() { t.g = T.golango }



