// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

extern void panic_callback();

void call_callback(void) {
	panic_callback();
}
