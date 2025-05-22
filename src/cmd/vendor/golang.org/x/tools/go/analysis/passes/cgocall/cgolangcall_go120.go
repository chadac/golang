// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !golang1.21

package cgolangcall

import "golang/types"

func setGoVersion(tc *types.Config, pkg *types.Package) {
	// no types.Package.GoVersion until Go 1.21
}
