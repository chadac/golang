// Copyright 2013 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gccgolang

#include "_cgolang_export.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

/* Test calling panic from C.  This is what SWIG does.  */

extern void _cgolang_panic(const char *);
extern void *_cgolang_allocate(size_t);

void
callPanic(void)
{
	_cgolang_panic("panic from C");
}
