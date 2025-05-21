// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build js && wasm

package os

// supportsCloseOnExec reports whether the platform supports the
// O_CLOEXEC flag.
const supportsCloseOnExec = false
