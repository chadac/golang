// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !plan9 && !wasm

package runtime

import "unsafe"

const isSbrkPlatform = false

func sysReserveAlignedSbrk(size, align uintptr) (unsafe.Pointer, uintptr) {
	panic("unreachable")
}
