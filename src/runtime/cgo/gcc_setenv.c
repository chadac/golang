// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix

#include "libcgolang.h"

#include <stdlib.h>

/* Stub for calling setenv */
void
x_cgolang_setenv(char **arg)
{
	_cgolang_tsan_acquire();
	setenv(arg[0], arg[1], 1);
	_cgolang_tsan_release();
}

/* Stub for calling unsetenv */
void
x_cgolang_unsetenv(char **arg)
{
	_cgolang_tsan_acquire();
	unsetenv(arg[0]);
	_cgolang_tsan_release();
}
