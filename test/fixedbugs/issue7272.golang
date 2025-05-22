// compile

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 7272: test builtin functions in statement context and in
// golang/defer functions.

package p

func F() {
	var a []int
	var c chan int
	var m map[int]int

	close(c)
	copy(a, a)
	delete(m, 0)
	panic(0)
	print("foo")
	println("bar")
	recover()

	(close(c))
	(copy(a, a))
	(delete(m, 0))
	(panic(0))
	(print("foo"))
	(println("bar"))
	(recover())

	golang close(c)
	golang copy(a, a)
	golang delete(m, 0)
	golang panic(0)
	golang print("foo")
	golang println("bar")
	golang recover()

	defer close(c)
	defer copy(a, a)
	defer delete(m, 0)
	defer panic(0)
	defer print("foo")
	defer println("bar")
	defer recover()
}
