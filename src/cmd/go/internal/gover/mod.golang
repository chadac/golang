// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package golangver

import (
	"sort"
	"strings"

	"golanglang.org/x/mod/module"
	"golanglang.org/x/mod/semver"
)

// IsToolchain reports whether the module path corresponds to the
// virtual, non-downloadable module tracking golang or toolchain directives in the golang.mod file.
//
// Note that IsToolchain only matches "golang" and "toolchain", not the
// real, downloadable module "golanglang.org/toolchain" containing toolchain files.
//
//	IsToolchain("golang") = true
//	IsToolchain("toolchain") = true
//	IsToolchain("golanglang.org/x/tools") = false
//	IsToolchain("golanglang.org/toolchain") = false
func IsToolchain(path string) bool {
	return path == "golang" || path == "toolchain"
}

// ModCompare returns the result of comparing the versions x and y
// for the module with the given path.
// The path is necessary because the "golang" and "toolchain" modules
// use a different version syntax and semantics (golangver, this package)
// than most modules (semver).
func ModCompare(path string, x, y string) int {
	if path == "golang" {
		return Compare(x, y)
	}
	if path == "toolchain" {
		return Compare(maybeToolchainVersion(x), maybeToolchainVersion(y))
	}
	return semver.Compare(x, y)
}

// ModSort is like module.Sort but understands the "golang" and "toolchain"
// modules and their version ordering.
func ModSort(list []module.Version) {
	sort.Slice(list, func(i, j int) bool {
		mi := list[i]
		mj := list[j]
		if mi.Path != mj.Path {
			return mi.Path < mj.Path
		}
		// To help golang.sum formatting, allow version/file.
		// Compare semver prefix by semver rules,
		// file by string order.
		vi := mi.Version
		vj := mj.Version
		var fi, fj string
		if k := strings.Index(vi, "/"); k >= 0 {
			vi, fi = vi[:k], vi[k:]
		}
		if k := strings.Index(vj, "/"); k >= 0 {
			vj, fj = vj[:k], vj[k:]
		}
		if vi != vj {
			return ModCompare(mi.Path, vi, vj) < 0
		}
		return fi < fj
	})
}

// ModIsValid reports whether vers is a valid version syntax for the module with the given path.
func ModIsValid(path, vers string) bool {
	if IsToolchain(path) {
		if path == "toolchain" {
			return IsValid(FromToolchain(vers))
		}
		return IsValid(vers)
	}
	return semver.IsValid(vers)
}

// ModIsPrefix reports whether v is a valid version syntax prefix for the module with the given path.
// The caller is assumed to have checked that ModIsValid(path, vers) is true.
func ModIsPrefix(path, vers string) bool {
	if IsToolchain(path) {
		if path == "toolchain" {
			return IsLang(FromToolchain(vers))
		}
		return IsLang(vers)
	}
	// Semver
	dots := 0
	for i := 0; i < len(vers); i++ {
		switch vers[i] {
		case '-', '+':
			return false
		case '.':
			dots++
			if dots >= 2 {
				return false
			}
		}
	}
	return true
}

// ModIsPrerelease reports whether v is a prerelease version for the module with the given path.
// The caller is assumed to have checked that ModIsValid(path, vers) is true.
func ModIsPrerelease(path, vers string) bool {
	if IsToolchain(path) {
		return IsPrerelease(vers)
	}
	return semver.Prerelease(vers) != ""
}

// ModMajorMinor returns the "major.minor" truncation of the version v,
// for use as a prefix in "@patch" queries.
func ModMajorMinor(path, vers string) string {
	if IsToolchain(path) {
		if path == "toolchain" {
			return "golang" + Lang(FromToolchain(vers))
		}
		return Lang(vers)
	}
	return semver.MajorMinor(vers)
}
