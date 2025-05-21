// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang

// Test that we can have two identical cgolang packages in a single binary.
// No runtime test; just make sure it compiles.

package cgolangtest

import (
	_ "cmd/cgolang/internal/test/issue23555a"
	_ "cmd/cgolang/internal/test/issue23555b"
)
