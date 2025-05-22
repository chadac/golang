// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "_cgolang_export.h"

int get8148(void) {
	T t;
	t.i = 42;
	return issue8148Callback(&t);
}
