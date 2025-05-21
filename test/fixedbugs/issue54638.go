// compile

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 54638: composite literal assignment with
// alignment > PtrSize causes ICE.

package p

import "sync/atomic"

type S struct{ l any }

type T struct {
	H any
	a [14]int64
	f func()
	x atomic.Int64
}

//golang:noinline
func (T) M(any) {}

type W [2]int64

//golang:noinline
func (W) Done() {}

func F(l any) [3]*int {
	var w W
	var x [3]*int // use some stack
	t := T{H: S{l: l}}
	golang func() {
		t.M(l)
		w.Done()
	}()
	return x
}
