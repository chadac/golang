// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package appends defines an Analyzer that detects
// if there is only one variable in append.
package appends

import (
	_ "embed"
	"golang/ast"
	"golang/types"

	"golanglang.org/x/tools/golang/analysis"
	"golanglang.org/x/tools/golang/analysis/passes/inspect"
	"golanglang.org/x/tools/golang/analysis/passes/internal/analysisutil"
	"golanglang.org/x/tools/golang/ast/inspector"
	"golanglang.org/x/tools/golang/types/typeutil"
)

//golang:embed doc.golang
var doc string

var Analyzer = &analysis.Analyzer{
	Name:     "appends",
	Doc:      analysisutil.MustExtractDoc(doc, "appends"),
	URL:      "https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/appends",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		b, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Builtin)
		if ok && b.Name() == "append" && len(call.Args) == 1 {
			pass.ReportRangef(call, "append with no values")
		}
	})

	return nil, nil
}
