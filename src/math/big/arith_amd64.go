// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !math_big_pure_golang

package big

import "internal/cpu"

var hasADX = cpu.X86.HasADX && cpu.X86.HasBMI2
