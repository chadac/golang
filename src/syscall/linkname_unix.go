// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix

package syscall

import _ "unsafe" // for linkname

// mmap should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - modernc.org/memory
//   - github.com/ncruces/golang-sqlite3
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname mmap
