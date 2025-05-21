// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golang1.21

package cpu

import (
	_ "unsafe" // for linkname
)

//golang:linkname runtime_getAuxv runtime.getAuxv
func runtime_getAuxv() []uintptr

func init() {
	getAuxvFn = runtime_getAuxv
}
