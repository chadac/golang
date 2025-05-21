// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package golangver

import "testing"

func TestFromToolchain(t *testing.T) { test1(t, fromToolchainTests, "FromToolchain", FromToolchain) }

var fromToolchainTests = []testCase1[string, string]{
	{"golang1.2.3", "1.2.3"},
	{"1.2.3", ""},
	{"golang1.2.3+bigcorp", ""},
	{"golang1.2.3-bigcorp", "1.2.3"},
	{"golang1.2.3-bigcorp more text", "1.2.3"},
	{"gccgolang-golang1.23rc4", ""},
	{"gccgolang-golang1.23rc4-bigdwarf", ""},
}
