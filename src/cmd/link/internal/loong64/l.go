// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package loong64

const (
	maxAlign  = 32 // max data alignment
	minAlign  = 1  // min data alignment
	funcAlign = 16
)

/* Used by ../../internal/ld/dwarf.golang */
const (
	dwarfRegSP = 3
	dwarfRegLR = 1
)
