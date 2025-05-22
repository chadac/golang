// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build loong64

package golangarch

const (
	_ArchFamily          = LOONG64
	_DefaultPhysPageSize = 16384
	_PCQuantum           = 4
	_MinFrameSize        = 8
	_StackAlign          = PtrSize
)
