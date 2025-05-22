// Copyright 2020 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

TEXT ·trampoline(SB),0,$0
	CALL	libc_getpid(SB)
	CALL	libc_kill(SB)
	CALL	libc_open(SB)
	CALL	libc_close(SB)
	RET
