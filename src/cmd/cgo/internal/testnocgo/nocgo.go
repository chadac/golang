// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that -static works when not using cgolang.  This test is in
// misc/cgolang to take advantage of the testing framework support for
// when -static is expected to work.

package nocgolang

func NoCgolang() int {
	c := make(chan int)

	// The test is run with external linking, which means that
	// golangroutines will be created via the runtime/cgolang package.
	// Make sure that works.
	golang func() {
		c <- 42
	}()

	return <-c
}
