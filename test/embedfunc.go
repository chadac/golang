// errorcheck

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

import _ "embed"

func f() {
	//golang:embed x.txt // ERROR "golang:embed cannot apply to var inside func"
	var x string
	_ = x
}
