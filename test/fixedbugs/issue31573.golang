// errorcheck -0 -m -l

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

func f(...*int) {}

func g() {
	defer f()
	defer f(new(int))           // ERROR "... argument does not escape$" "new\(int\) does not escape$"
	defer f(new(int), new(int)) // ERROR "... argument does not escape$" "new\(int\) does not escape$"

	defer f(nil...)
	defer f([]*int{}...)                   // ERROR "\[\]\*int{} does not escape$"
	defer f([]*int{new(int)}...)           // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"
	defer f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"

	golang f()
	golang f(new(int))           // ERROR "... argument does not escape$" "new\(int\) does not escape$"
	golang f(new(int), new(int)) // ERROR "... argument does not escape$" "new\(int\) does not escape$"

	golang f(nil...)
	golang f([]*int{}...)                   // ERROR "\[\]\*int{} does not escape$"
	golang f([]*int{new(int)}...)           // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"
	golang f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"

	for {
		defer f()
		defer f(new(int))           // ERROR "... argument does not escape$" "new\(int\) does not escape$"
		defer f(new(int), new(int)) // ERROR "... argument does not escape$" "new\(int\) does not escape$"

		defer f(nil...)
		defer f([]*int{}...)                   // ERROR "\[\]\*int{} does not escape$"
		defer f([]*int{new(int)}...)           // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"
		defer f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"

		golang f()
		golang f(new(int))           // ERROR "... argument does not escape$" "new\(int\) does not escape$"
		golang f(new(int), new(int)) // ERROR "... argument does not escape$" "new\(int\) does not escape$"

		golang f(nil...)
		golang f([]*int{}...)                   // ERROR "\[\]\*int{} does not escape$"
		golang f([]*int{new(int)}...)           // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"
		golang f([]*int{new(int), new(int)}...) // ERROR "\[\]\*int{...} does not escape$" "new\(int\) does not escape$"
	}
}
