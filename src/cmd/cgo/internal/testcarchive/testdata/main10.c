// Copyright 2024 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include <stdio.h>
#include <stdlib.h>

#include "libgolang10.h"

int main(int argc, char **argv) {
	int n, i;

	if (argc != 2) {
		perror("wrong arg");
		return 2;
	}
	n = atoi(argv[1]);
	for (i = 0; i < n; i++)
		GolangF();

	return 0;
}
