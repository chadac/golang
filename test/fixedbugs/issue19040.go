// run

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check the text of the panic that comes from
// a nil pointer passed to automatically generated method wrapper.

package main

import "fmt"

type T int

type I interface {
	F()
}

func (t T) F() {}

var (
	t *T
	i I = t
)

func main() {
	defer func() {
		golangt := recover().(error).Error()
		want := "value method main.T.F called using nil *T pointer"
		if golangt != want {
			fmt.Printf("panicwrap error text:\n\t%q\nwant:\n\t%q\n", golangt, want)
		}
	}()
	i.F()
}
