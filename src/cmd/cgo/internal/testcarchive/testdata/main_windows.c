// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Dummy implementations for Windows, because Windows doesn't
 * support Unix-style signal handling.
 */

int install_handler() {
	return 0;
}


int check_handler() {
	return 0;
}
