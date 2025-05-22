// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package gc

import "internal/golangarch"

// ObjMask is a bitmap where each bit corresponds to an object in a span.
//
// It is sized to accomodate all size classes.
type ObjMask [MaxObjsPerSpan / (golangarch.PtrSize * 8)]uintptr

// PtrMask is a bitmap where each bit represents a pointer-word in a single runtime page.
type PtrMask [PageSize / golangarch.PtrSize / (golangarch.PtrSize * 8)]uintptr
