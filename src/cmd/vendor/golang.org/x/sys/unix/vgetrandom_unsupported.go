// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !linux || !golang1.24

package unix

func vgetrandom(p []byte, flags uint32) (ret int, supported bool) {
	return -1, false
}
