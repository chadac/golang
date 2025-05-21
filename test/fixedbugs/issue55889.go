// errorcheck -0 -lang=golang1.17

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Prior to Go 1.18, ineffectual //golang:linkname directives were treated
// as noops. Ensure that modules that contain these directives (e.g.,
// x/sys prior to golang.dev/cl/274573) continue to compile.

package p

import _ "unsafe"

//golang:linkname nonexistent nonexistent

//golang:linkname constant constant
const constant = 42

//golang:linkname typename typename
type typename int
