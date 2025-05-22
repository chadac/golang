// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gccgolang

package build

import "runtime"

// getToolDir returns the default value of ToolDir.
func getToolDir() string {
	return envOr("GCCGOTOOLDIR", runtime.GCCGOTOOLDIR)
}
