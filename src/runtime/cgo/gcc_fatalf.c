// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || (!android && linux) || dragolangnfly || freebsd || netbsd || openbsd || solaris

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include "libcgolang.h"

void
fatalf(const char* format, ...)
{
	va_list ap;

	fprintf(stderr, "runtime/cgolang: ");
	va_start(ap, format);
	vfprintf(stderr, format, ap);
	va_end(ap);
	fprintf(stderr, "\n");
	abort();
}
