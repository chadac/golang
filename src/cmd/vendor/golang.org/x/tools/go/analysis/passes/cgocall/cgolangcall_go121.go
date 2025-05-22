// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golang1.21

package cgolangcall

import "golang/types"

func setGolangVersion(tc *types.Config, pkg *types.Package) {
	tc.GolangVersion = pkg.GolangVersion()
}
