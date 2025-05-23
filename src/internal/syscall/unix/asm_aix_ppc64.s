// Copyright 2018 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

//
// System calls for aix/ppc64 are implemented in syscall/syscall_aix.golang
//

TEXT ·syscall6(SB),NOSPLIT,$0
	JMP	syscall·syscall6(SB)
