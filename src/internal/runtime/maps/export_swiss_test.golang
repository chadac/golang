// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golangexperiment.swissmap

package maps

import (
	"internal/abi"
	"unsafe"
)

func newTestMapType[K comparable, V any]() *abi.SwissMapType {
	var m map[K]V
	mTyp := abi.TypeOf(m)
	mt := (*abi.SwissMapType)(unsafe.Pointer(mTyp))
	return mt
}
