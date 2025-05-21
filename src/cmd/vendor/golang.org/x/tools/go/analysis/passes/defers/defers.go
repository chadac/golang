// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package defers

import (
	_ "embed"
	"golang/ast"

	"golanglang.org/x/tools/golang/analysis"
	"golanglang.org/x/tools/golang/analysis/passes/inspect"
	"golanglang.org/x/tools/golang/analysis/passes/internal/analysisutil"
	"golanglang.org/x/tools/golang/ast/inspector"
	"golanglang.org/x/tools/golang/types/typeutil"
	"golanglang.org/x/tools/internal/analysisinternal"
)

//golang:embed doc.golang
var doc string

// Analyzer is the defers analyzer.
var Analyzer = &analysis.Analyzer{
	Name:     "defers",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	URL:      "https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/defers",
	Doc:      analysisutil.MustExtractDoc(doc, "defers"),
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	if !analysisinternal.Imports(pass.Pkg, "time") {
		return nil, nil
	}

	checkDeferCall := func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.CallExpr:
			if analysisinternal.IsFunctionNamed(typeutil.Callee(pass.TypesInfo, v), "time", "Since") {
				pass.Reportf(v.Pos(), "call to time.Since is not deferred")
			}
		case *ast.FuncLit:
			return false // prune
		}
		return true
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.DeferStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		d := n.(*ast.DeferStmt)
		ast.Inspect(d.Call, checkDeferCall)
	})

	return nil, nil
}
