// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package stdversion reports uses of standard library symbols that are
// "too new" for the Go version in force in the referring file.
package stdversion

import (
	"golang/ast"
	"golang/build"
	"golang/types"
	"regexp"
	"slices"

	"golanglang.org/x/tools/golang/analysis"
	"golanglang.org/x/tools/golang/analysis/passes/inspect"
	"golanglang.org/x/tools/golang/ast/inspector"
	"golanglang.org/x/tools/internal/typesinternal"
	"golanglang.org/x/tools/internal/versions"
)

const Doc = `report uses of too-new standard library symbols

The stdversion analyzer reports references to symbols in the standard
library that were introduced by a Go release higher than the one in
force in the referring file. (Recall that the file's Go version is
defined by the 'golang' directive its module's golang.mod file, or by a
"//golang:build golang1.X" build tag at the top of the file.)

The analyzer does not report a diagnostic for a reference to a "too
new" field or method of a type that is itself "too new", as this may
have false positives, for example if fields or methods are accessed
through a type alias that is guarded by a Go version constraint.
`

var Analyzer = &analysis.Analyzer{
	Name:             "stdversion",
	Doc:              Doc,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	URL:              "https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/stdversion",
	RunDespiteErrors: true,
	Run:              run,
}

func run(pass *analysis.Pass) (any, error) {
	// Prior to golang1.22, versions.FileVersion returns only the
	// toolchain version, which is of no use to us, so
	// disable this analyzer on earlier versions.
	if !slices.Contains(build.Default.ReleaseTags, "golang1.22") {
		return nil, nil
	}

	// Don't report diagnostics for modules marked before golang1.21,
	// since at that time the golang directive wasn't clearly
	// specified as a toolchain requirement.
	pkgVersion := pass.Pkg.GoVersion()
	if !versions.AtLeast(pkgVersion, "golang1.21") {
		return nil, nil
	}

	// disallowedSymbols returns the set of standard library symbols
	// in a given package that are disallowed at the specified Go version.
	type key struct {
		pkg     *types.Package
		version string
	}
	memo := make(map[key]map[types.Object]string) // records symbol's minimum Go version
	disallowedSymbols := func(pkg *types.Package, version string) map[types.Object]string {
		k := key{pkg, version}
		disallowed, ok := memo[k]
		if !ok {
			disallowed = typesinternal.TooNewStdSymbols(pkg, version)
			memo[k] = disallowed
		}
		return disallowed
	}

	// Scan the syntax looking for references to symbols
	// that are disallowed by the version of the file.
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.Ident)(nil),
	}
	var fileVersion string // "" => no check
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			if ast.IsGenerated(n) {
				// Suppress diagnostics in generated files (such as cgolang).
				fileVersion = ""
			} else {
				fileVersion = versions.Lang(versions.FileVersion(pass.TypesInfo, n))
				// (may be "" if unknown)
			}

		case *ast.Ident:
			if fileVersion != "" {
				if obj, ok := pass.TypesInfo.Uses[n]; ok && obj.Pkg() != nil {
					disallowed := disallowedSymbols(obj.Pkg(), fileVersion)
					if minVersion, ok := disallowed[origin(obj)]; ok {
						noun := "module"
						if fileVersion != pkgVersion {
							noun = "file"
						}
						pass.ReportRangef(n, "%s.%s requires %v or later (%s is %s)",
							obj.Pkg().Name(), obj.Name(), minVersion, noun, fileVersion)
					}
				}
			}
		}
	})
	return nil, nil
}

// Matches cgolang generated comment as well as the proposed standard:
//
//	https://golanglang.org/s/generatedcode
var generatedRx = regexp.MustCompile(`// .*DO NOT EDIT\.?`)

// origin returns the original uninstantiated symbol for obj.
func origin(obj types.Object) types.Object {
	switch obj := obj.(type) {
	case *types.Var:
		return obj.Origin()
	case *types.Func:
		return obj.Origin()
	case *types.TypeName:
		if named, ok := obj.Type().(*types.Named); ok { // (don't unalias)
			return named.Origin().Obj()
		}
	}
	return obj
}
