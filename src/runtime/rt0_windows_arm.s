// Copyright 2018 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "golang_asm.h"
#include "golang_tls.h"
#include "textflag.h"

// This is the entry point for the program from the
// kernel for an ordinary -buildmode=exe program.
TEXT _rt0_arm_windows(SB),NOSPLIT|NOFRAME,$0
	B	·rt0_golang(SB)
