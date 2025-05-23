// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package golangver implements support for Golang toolchain versions like 1.21.0 and 1.21rc1.
// (For historical reasons, Golang does not use semver for its toolchains.)
// This package provides the same basic analysis that golanglang.org/x/mod/semver does for semver.
// It also provides some helpers for extracting versions from golang.mod files
// and for dealing with module.Versions that may use Golang versions or semver
// depending on the module path.
package golangver

import (
	"internal/golangver"
)

// Compare returns -1, 0, or +1 depending on whether
// x < y, x == y, or x > y, interpreted as toolchain versions.
// The versions x and y must not begin with a "golang" prefix: just "1.21" not "golang1.21".
// Malformed versions compare less than well-formed versions and equal to each other.
// The language version "1.21" compares less than the release candidate and eventual releases "1.21rc1" and "1.21.0".
func Compare(x, y string) int {
	return golangver.Compare(x, y)
}

// Max returns the maximum of x and y interpreted as toolchain versions,
// compared using Compare.
// If x and y compare equal, Max returns x.
func Max(x, y string) string {
	return golangver.Max(x, y)
}

// IsLang reports whether v denotes the overall Golang language version
// and not a specific release. Starting with the Golang 1.21 release, "1.x" denotes
// the overall language version; the first release is "1.x.0".
// The distinction is important because the relative ordering is
//
//	1.21 < 1.21rc1 < 1.21.0
//
// meaning that Golang 1.21rc1 and Golang 1.21.0 will both handle golang.mod files that
// say "golang 1.21", but Golang 1.21rc1 will not handle files that say "golang 1.21.0".
func IsLang(x string) bool {
	return golangver.IsLang(x)
}

// Lang returns the Golang language version. For example, Lang("1.2.3") == "1.2".
func Lang(x string) string {
	return golangver.Lang(x)
}

// IsPrerelease reports whether v denotes a Golang prerelease version.
func IsPrerelease(x string) bool {
	return golangver.Parse(x).Kind != ""
}

// Prev returns the Golang major release immediately preceding v,
// or v itself if v is the first Golang major release (1.0) or not a supported
// Golang version.
//
// Examples:
//
//	Prev("1.2") = "1.1"
//	Prev("1.3rc4") = "1.2"
func Prev(x string) string {
	v := golangver.Parse(x)
	if golangver.CmpInt(v.Minor, "1") <= 0 {
		return v.Major
	}
	return v.Major + "." + golangver.DecInt(v.Minor)
}

// IsValid reports whether the version x is valid.
func IsValid(x string) bool {
	return golangver.IsValid(x)
}
