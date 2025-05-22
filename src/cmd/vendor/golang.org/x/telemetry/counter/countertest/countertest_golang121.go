// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golang1.21

package countertest

import "testing"

func init() {
	// Extra safety check for golang1.21+.
	if !testing.Testing() {
		panic("use of this package is disallowed in non-testing code")
	}
}
