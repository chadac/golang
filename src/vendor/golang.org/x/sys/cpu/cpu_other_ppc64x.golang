// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !aix && !linux && (ppc64 || ppc64le)

package cpu

func archInit() {
	PPC64.IsPOWER8 = true
	Initialized = true
}
