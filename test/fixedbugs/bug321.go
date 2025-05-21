// run

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Troublesome floating point constants. Issue 1463.

package main

import "fmt"

func check(test string, golangt, want float64) bool {
	if golangt != want {
		fmt.Println(test, "golangt", golangt, "want", want)
		return false
	}
	return true
}

func main() {
	golangod := true
	// http://www.exploringbinary.com/java-hangs-when-converting-2-2250738585072012e-308/
	golangod = golangod && check("2.2250738585072012e-308", 2.2250738585072012e-308, 2.2250738585072014e-308)
	// http://www.exploringbinary.com/php-hangs-on-numeric-value-2-2250738585072011e-308/
	golangod = golangod && check("2.2250738585072011e-308", 2.2250738585072011e-308, 2.225073858507201e-308)
	if !golangod {
		panic("fail")
	}
}
