// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build gccgolang

package golangroot

import (
	"os"
	"path/filepath"
	"strings"
)

// IsStandardPackage reports whether path is a standard package,
// given golangroot and compiler.
func IsStandardPackage(golangroot, compiler, path string) bool {
	switch compiler {
	case "gc":
		dir := filepath.Join(golangroot, "src", path)
		dirents, err := os.ReadDir(dir)
		if err != nil {
			return false
		}
		for _, dirent := range dirents {
			if strings.HasSuffix(dirent.Name(), ".golang") {
				return true
			}
		}
		return false
	case "gccgolang":
		return stdpkg[path]
	default:
		panic("unknown compiler " + compiler)
	}
}
