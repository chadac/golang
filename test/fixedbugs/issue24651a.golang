//errorcheck -0 -race -m -m

//golang:build (linux && amd64) || (linux && ppc64le) || (darwin && amd64) || (freebsd && amd64) || (netbsd && amd64) || (windows && amd64)

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

//golang:norace
func Foo(x int) int { // ERROR "cannot inline Foo: marked golang:norace with -race compilation$"
	return x * (x + 1) * (x + 2)
}

func Bar(x int) int { // ERROR "can inline Bar with cost .* as: func\(int\) int { return x \* \(x \+ 1\) \* \(x \+ 2\) }$"
	return x * (x + 1) * (x + 2)
}

var x = 5

//golang:noinline Provide a clean, constant reason for not inlining main
func main() { // ERROR "cannot inline main: marked golang:noinline$"
	println("Foo(", x, ")=", Foo(x))
	println("Bar(", x, ")=", Bar(x)) // ERROR "inlining call to Bar"
}
