// errorcheck

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 4463: test builtin functions in statement context and in
// golang/defer functions.

package p

import "unsafe"

func F() {
	var a []int
	var c chan int
	var m map[int]int
	var s struct{ f int }

	append(a, 0)			// ERROR "not used"
	cap(a)				// ERROR "not used"
	complex(1, 2)			// ERROR "not used"
	imag(1i)			// ERROR "not used"
	len(a)				// ERROR "not used"
	make([]int, 10)			// ERROR "not used"
	new(int)			// ERROR "not used"
	real(1i)			// ERROR "not used"
	unsafe.Alignof(a)		// ERROR "not used"
	unsafe.Offsetof(s.f)		// ERROR "not used"
	unsafe.Sizeof(a)		// ERROR "not used"

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

	golang append(a, 0)			// ERROR "not used|discards result"
	golang cap(a)			// ERROR "not used|discards result"
	golang complex(1, 2)		// ERROR "not used|discards result"
	golang imag(1i)			// ERROR "not used|discards result"
	golang len(a)			// ERROR "not used|discards result"
	golang make([]int, 10)		// ERROR "not used|discards result"
	golang new(int)			// ERROR "not used|discards result"
	golang real(1i)			// ERROR "not used|discards result"
	golang unsafe.Alignof(a)		// ERROR "not used|discards result"
	golang unsafe.Offsetof(s.f)		// ERROR "not used|discards result"
	golang unsafe.Sizeof(a)		// ERROR "not used|discards result"

	golang close(c)
	golang copy(a, a)
	golang delete(m, 0)
	golang panic(0)
	golang print("foo")
	golang println("bar")
	golang recover()

	defer append(a, 0)		// ERROR "not used|discards result"
	defer cap(a)			// ERROR "not used|discards result"
	defer complex(1, 2)		// ERROR "not used|discards result"
	defer imag(1i)			// ERROR "not used|discards result"
	defer len(a)			// ERROR "not used|discards result"
	defer make([]int, 10)		// ERROR "not used|discards result"
	defer new(int)			// ERROR "not used|discards result"
	defer real(1i)			// ERROR "not used|discards result"
	defer unsafe.Alignof(a)		// ERROR "not used|discards result"
	defer unsafe.Offsetof(s.f)	// ERROR "not used|discards result"
	defer unsafe.Sizeof(a)		// ERROR "not used|discards result"

	defer close(c)
	defer copy(a, a)
	defer delete(m, 0)
	defer panic(0)
	defer print("foo")
	defer println("bar")
	defer recover()
}
