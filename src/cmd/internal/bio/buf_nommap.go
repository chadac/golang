// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !unix

package bio

func (r *Reader) sliceOS(length uint64) ([]byte, bool) {
	return nil, false
}
