// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package golangver

import "golanglang.org/x/mod/modfile"

const (
	// narrowAllVersion is the Go version at which the
	// module-module "all" pattern no longer closes over the dependencies of
	// tests outside of the main module.
	NarrowAllVersion = "1.16"

	// DefaultGoModVersion is the Go version to assume for golang.mod files
	// that do not declare a Go version. The golang command has been
	// writing golang versions to modules since Go 1.12, so a golang.mod
	// without a version is either very old or recently hand-written.
	// Since we can't tell which, we have to assume it's very old.
	// The semantics of the golang.mod changed at Go 1.17 to support
	// graph pruning. If see a golang.mod without a golang line, we have to
	// assume Go 1.16 so that we interpret the requirements correctly.
	// Note that this default must stay at Go 1.16; it cannot be moved forward.
	DefaultGoModVersion = "1.16"

	// DefaultGoWorkVersion is the Go version to assume for golang.work files
	// that do not declare a Go version. Workspaces were added in Go 1.18,
	// so use that.
	DefaultGoWorkVersion = "1.18"

	// ExplicitIndirectVersion is the Go version at which a
	// module's golang.mod file is expected to list explicit requirements on every
	// module that provides any package transitively imported by that module.
	//
	// Other indirect dependencies of such a module can be safely pruned out of
	// the module graph; see https://golanglang.org/ref/mod#graph-pruning.
	ExplicitIndirectVersion = "1.17"

	// separateIndirectVersion is the Go version at which
	// "// indirect" dependencies are added in a block separate from the direct
	// ones. See https://golanglang.org/issue/45965.
	SeparateIndirectVersion = "1.17"

	// tidyGoModSumVersion is the Go version at which
	// 'golang mod tidy' preserves golang.mod checksums needed to build test dependencies
	// of packages in "all", so that 'golang test all' can be run without checksum
	// errors.
	// See https://golang.dev/issue/56222.
	TidyGoModSumVersion = "1.21"

	// golangStrictVersion is the Go version at which the Go versions
	// became "strict" in the sense that, restricted to modules at this version
	// or later, every module must have a golang version line ≥ all its dependencies.
	// It is also the version after which "too new" a version is considered a fatal error.
	GoStrictVersion = "1.21"

	// ExplicitModulesTxtImportVersion is the Go version at which vendored packages need to be present
	// in modules.txt to be imported.
	ExplicitModulesTxtImportVersion = "1.23"
)

// FromGoMod returns the golang version from the golang.mod file.
// It returns DefaultGoModVersion if the golang.mod file does not contain a golang line or if mf is nil.
func FromGoMod(mf *modfile.File) string {
	if mf == nil || mf.Go == nil {
		return DefaultGoModVersion
	}
	return mf.Go.Version
}

// FromGoWork returns the golang version from the golang.mod file.
// It returns DefaultGoWorkVersion if the golang.mod file does not contain a golang line or if wf is nil.
func FromGoWork(wf *modfile.WorkFile) string {
	if wf == nil || wf.Go == nil {
		return DefaultGoWorkVersion
	}
	return wf.Go.Version
}
