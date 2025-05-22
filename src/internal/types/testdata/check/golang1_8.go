// -lang=golang1.8

// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Check Golang language version-specific errors.

package p

// type alias declarations
type any = /* ERROR "type alias requires golang1.9 or later" */ interface{}
