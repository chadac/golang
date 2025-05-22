// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang

package cgolangtest

import (
	"testing"

	"cmd/cgolang/internal/test/gcc68255"
)

func testGCC68255(t *testing.T) {
	if !gcc68255.F() {
		t.Error("C global variable was not initialized")
	}
}
