// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build arm64 || ppc64 || ppc64le

package math

const haveArchModf = true

func archModf(f float64) (int float64, frac float64)
