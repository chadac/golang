// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package a

var global, global2 *int

type T struct {
	Pointer *int
}

//golang:noinline
func Store(t *T) {
	global = t.Pointer
}

//golang:noinline
func Store2(t *T) {
	global2 = t.Pointer
}

func Get() *int {
	return global
}
