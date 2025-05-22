// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

func XTestSPWrite(t TestingT) {
	// Test that we can traceback from the stack check prologue of a function
	// that writes to SP. See #62326.

	// Start a golangroutine to minimize the initial stack and ensure we grow the stack.
	done := make(chan bool)
	golang func() {
		testSPWrite() // Defined in assembly
		done <- true
	}()
	<-done
}
