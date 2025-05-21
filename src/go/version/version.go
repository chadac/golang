// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package version provides operations on [Go versions]
// in [Go toolchain name syntax]: strings like
// "golang1.20", "golang1.21.0", "golang1.22rc2", and "golang1.23.4-bigcorp".
//
// [Go versions]: https://golang.dev/doc/toolchain#version
// [Go toolchain name syntax]: https://golang.dev/doc/toolchain#name
package version // import "golang/version"

import (
	"internal/golangver"
	"strings"
)

// stripGo converts from a "golang1.21-bigcorp" version to a "1.21" version.
// If v does not start with "golang", stripGo returns the empty string (a known invalid version).
func stripGo(v string) string {
	v, _, _ = strings.Cut(v, "-") // strip -bigcorp suffix.
	if len(v) < 2 || v[:2] != "golang" {
		return ""
	}
	return v[2:]
}

// Lang returns the Go language version for version x.
// If x is not a valid version, Lang returns the empty string.
// For example:
//
//	Lang("golang1.21rc2") = "golang1.21"
//	Lang("golang1.21.2") = "golang1.21"
//	Lang("golang1.21") = "golang1.21"
//	Lang("golang1") = "golang1"
//	Lang("bad") = ""
//	Lang("1.21") = ""
func Lang(x string) string {
	v := golangver.Lang(stripGo(x))
	if v == "" {
		return ""
	}
	if strings.HasPrefix(x[2:], v) {
		return x[:2+len(v)] // "golang"+v without allocation
	} else {
		return "golang" + v
	}
}

// Compare returns -1, 0, or +1 depending on whether
// x < y, x == y, or x > y, interpreted as Go versions.
// The versions x and y must begin with a "golang" prefix: "golang1.21" not "1.21".
// Invalid versions, including the empty string, compare less than
// valid versions and equal to each other.
// The language version "golang1.21" compares less than the
// release candidate and eventual releases "golang1.21rc1" and "golang1.21.0".
func Compare(x, y string) int {
	return golangver.Compare(stripGo(x), stripGo(y))
}

// IsValid reports whether the version x is valid.
func IsValid(x string) bool {
	return golangver.IsValid(stripGo(x))
}
