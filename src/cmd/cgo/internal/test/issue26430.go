// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang

// Issue 26430: incomplete typedef leads to inconsistent typedefs error.
// No runtime test; just make sure it compiles.

package cgolangtest

import _ "cmd/cgolang/internal/test/issue26430"
