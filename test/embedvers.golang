// errorcheck -lang=golang1.15

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

import _ "embed"

//golang:embed x.txt // ERROR "golang:embed requires golang1.16 or later"
var x string
