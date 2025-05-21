// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gc

#include "textflag.h"

//
// System calls for ppc64, AIX are implemented in runtime/syscall_aix.golang
//

TEXT ·syscall6(SB),NOSPLIT,$0-88
	JMP	syscall·syscall6(SB)

TEXT ·rawSyscall6(SB),NOSPLIT,$0-88
	JMP	syscall·rawSyscall6(SB)
