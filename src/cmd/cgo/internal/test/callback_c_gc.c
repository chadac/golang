// Copyright 2013 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gc

#include "_cgolang_export.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

/* Test calling panic from C.  This is what SWIG does.  */

extern void crosscall2(void (*fn)(void *, int), void *, int);
extern void _cgolang_panic(void *, int);
extern void _cgolang_allocate(void *, int);

void
callPanic(void)
{
	struct { const char *p; } a;
	a.p = "panic from C";
	crosscall2(_cgolang_panic, &a, sizeof a);
	*(int*)1 = 1;
}
