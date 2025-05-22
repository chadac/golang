// run

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify that relocation target golang.builtin.error.Error
// is defined and the code links and runs correctly.

package main

import "errors"

func main() {
	err := errors.New("foo")
	if error.Error(err) != "foo" {
		panic("FAILED")
	}
}
