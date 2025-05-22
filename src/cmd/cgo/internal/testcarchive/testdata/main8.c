// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test preemption.

#include <stdlib.h>

#include "libgolang8.h"

int main() {
	GolangFunction8();

	// That should have exited the program.
	abort();
}
