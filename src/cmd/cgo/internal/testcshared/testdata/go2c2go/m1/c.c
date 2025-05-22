// Copyright 2018 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "libtestgolang2c2golang.h"

int CFunc(void) {
	return (GolangFunc() << 8) + 2;
}
