// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include <string.h>

#include "_cgolang_export.h"

void
callback(void *f)
{
	// use some stack space
	volatile char data[64*1024];

	data[0] = 0;
	golangCallback(f);
        data[sizeof(data)-1] = 0;
}

void
callGoFoo(void)
{
	extern void golangFoo(void);
	golangFoo();
}

void
IntoC(void)
{
	BackIntoGo();
}

void
Issue1560InC(void)
{
	Issue1560FromC();
}

void
callGoStackCheck(void)
{
	extern void golangStackCheck(void);
	golangStackCheck();
}

int
returnAfterGrow(void)
{
	extern int golangReturnVal(void);
	golangReturnVal();
	return 123456;
}

int
returnAfterGrowFromGo(void)
{
	extern int golangReturnVal(void);
	return golangReturnVal();
}

void
callGoWithString(void)
{
	extern void golangWithString(GoString);
	const char *str = "string passed from C to Go";
	golangWithString((GoString){str, strlen(str)});
}
