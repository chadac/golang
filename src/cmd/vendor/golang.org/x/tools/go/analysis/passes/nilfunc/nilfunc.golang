// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package nilfunc defines an Analyzer that checks for useless
// comparisons against nil.
package nilfunc

import (
	_ "embed"
	"golang/ast"
	"golang/token"
	"golang/types"

	"golanglang.org/x/tools/golang/analysis"
	"golanglang.org/x/tools/golang/analysis/passes/inspect"
	"golanglang.org/x/tools/golang/analysis/passes/internal/analysisutil"
	"golanglang.org/x/tools/golang/ast/inspector"
	"golanglang.org/x/tools/internal/typesinternal"
)

//golang:embed doc.golang
var doc string

var Analyzer = &analysis.Analyzer{
	Name:     "nilfunc",
	Doc:      analysisutil.MustExtractDoc(doc, "nilfunc"),
	URL:      "https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/nilfunc",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.BinaryExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.BinaryExpr)

		// Only want == or != comparisons.
		if e.Op != token.EQL && e.Op != token.NEQ {
			return
		}

		// Only want comparisons with a nil identifier on one side.
		var e2 ast.Expr
		switch {
		case pass.TypesInfo.Types[e.X].IsNil():
			e2 = e.Y
		case pass.TypesInfo.Types[e.Y].IsNil():
			e2 = e.X
		default:
			return
		}

		// Only want functions.
		obj := pass.TypesInfo.Uses[typesinternal.UsedIdent(pass.TypesInfo, e2)]
		if _, ok := obj.(*types.Func); !ok {
			return
		}

		pass.ReportRangef(e, "comparison of function %v %v nil is always %v", obj.Name(), e.Op, e.Op == token.NEQ)
	})
	return nil, nil
}
