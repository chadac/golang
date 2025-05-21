// -lang=golang1.16

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golang1.21

package main

import "slices"

func main() {
	_ = slices.Clone([]string{}) // no error should be reported here
}
