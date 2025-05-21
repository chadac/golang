// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package cgolangtest

import "testing"

// This test actually doesn't have anything to do with cgolang.  It is a
// test of https://golanglang.org/issue/7234, a compiler/linker bug in
// handling string constants when using -linkmode=external.  The test
// is in this directory because we routinely test -linkmode=external
// here.

var v7234 = [...]string{"runtime/cgolang"}

func Test7234(t *testing.T) {
	if v7234[0] != "runtime/cgolang" {
		t.Errorf("bad string constant %q", v7234[0])
	}
}
