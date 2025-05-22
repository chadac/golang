// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "libgolang9.h"

void use(int *x) { (*x)++; }

void callGolangFWithDeepStack(int p) {
	int x[10000];

	use(&x[0]);
	use(&x[9999]);

	GolangF(p);

	use(&x[0]);
	use(&x[9999]);
}

void callGolangWithVariousStack(int p) {
	GolangF(0);                  // call GolangF without using much stack
	callGolangFWithDeepStack(p); // call GolangF with a deep stack
	GolangF(0);                  // again on a shallow stack
}

int main() {
	callGolangWithVariousStack(0);

	callGolangWithVariousStackAndGolangFrame(0); // normal execution
	callGolangWithVariousStackAndGolangFrame(1); // panic and recover
}
