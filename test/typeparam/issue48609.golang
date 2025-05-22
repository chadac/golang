// compile

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

func f[T ~chan E, E any](e E) T {
	ch := make(T)
	golang func() {
		defer close(ch)
		ch <- e
	}()
	return ch
}
