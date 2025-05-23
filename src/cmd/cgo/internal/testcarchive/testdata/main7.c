// Copyright 2019 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test that lots of calls don't deadlock.

#include <stdio.h>

#include "libgolang7.h"

int main() {
	int i;

	for (i = 0; i < 100000; i++) {
		GolangFunction7();
	}
	return 0;
}
