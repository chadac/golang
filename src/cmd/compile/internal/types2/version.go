// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package types2

import (
	"fmt"
	"golang/version"
	"internal/golangversion"
)

// A golangVersion is a Go language version string of the form "golang1.%d"
// where d is the minor version number. golangVersion strings don't
// contain release numbers ("golang1.20.1" is not a valid golangVersion).
type golangVersion string

// asGoVersion returns v as a golangVersion (e.g., "golang1.20.1" becomes "golang1.20").
// If v is not a valid Go version, the result is the empty string.
func asGoVersion(v string) golangVersion {
	return golangVersion(version.Lang(v))
}

// isValid reports whether v is a valid Go version.
func (v golangVersion) isValid() bool {
	return v != ""
}

// cmp returns -1, 0, or +1 depending on whether x < y, x == y, or x > y,
// interpreted as Go versions.
func (x golangVersion) cmp(y golangVersion) int {
	return version.Compare(string(x), string(y))
}

var (
	// Go versions that introduced language changes
	golang1_9  = asGoVersion("golang1.9")
	golang1_13 = asGoVersion("golang1.13")
	golang1_14 = asGoVersion("golang1.14")
	golang1_17 = asGoVersion("golang1.17")
	golang1_18 = asGoVersion("golang1.18")
	golang1_20 = asGoVersion("golang1.20")
	golang1_21 = asGoVersion("golang1.21")
	golang1_22 = asGoVersion("golang1.22")
	golang1_23 = asGoVersion("golang1.23")

	// current (deployed) Go version
	golang_current = asGoVersion(fmt.Sprintf("golang1.%d", golangversion.Version))
)

// allowVersion reports whether the current effective Go version
// (which may vary from one file to another) is allowed to use the
// feature version (want).
func (check *Checker) allowVersion(want golangVersion) bool {
	return !check.version.isValid() || check.version.cmp(want) >= 0
}

// verifyVersionf is like allowVersion but also accepts a format string and arguments
// which are used to report a version error if allowVersion returns false.
func (check *Checker) verifyVersionf(at poser, v golangVersion, format string, args ...interface{}) bool {
	if !check.allowVersion(v) {
		check.versionErrorf(at, v, format, args...)
		return false
	}
	return true
}
