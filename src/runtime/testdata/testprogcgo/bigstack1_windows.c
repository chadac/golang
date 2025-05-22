// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This is not in bigstack_windows.c because it needs to be part of
// testprogcgolang but is not part of the DLL built from bigstack_windows.c.

#include "_cgolang_export.h"

void CallGolangBigStack1(char* p) {
	golangBigStack1(p);
}
