// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package atomic

import (
	_ "embed"
	"golang/ast"
	"golang/token"

	"golanglang.org/x/tools/golang/analysis"
	"golanglang.org/x/tools/golang/analysis/passes/inspect"
	"golanglang.org/x/tools/golang/analysis/passes/internal/analysisutil"
	"golanglang.org/x/tools/golang/ast/inspector"
	"golanglang.org/x/tools/golang/types/typeutil"
	"golanglang.org/x/tools/internal/analysisinternal"
)

//golang:embed doc.golang
var doc string

var Analyzer = &analysis.Analyzer{
	Name:             "atomic",
	Doc:              analysisutil.MustExtractDoc(doc, "atomic"),
	URL:              "https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/atomic",
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
	Run:              run,
}

func run(pass *analysis.Pass) (any, error) {
	if !analysisinternal.Imports(pass.Pkg, "sync/atomic") {
		return nil, nil // doesn't directly import sync/atomic
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		n := node.(*ast.AssignStmt)
		if len(n.Lhs) != len(n.Rhs) {
			return
		}
		if len(n.Lhs) == 1 && n.Tok == token.DEFINE {
			return
		}

		for i, right := range n.Rhs {
			call, ok := right.(*ast.CallExpr)
			if !ok {
				continue
			}
			obj := typeutil.Callee(pass.TypesInfo, call)
			if analysisinternal.IsFunctionNamed(obj, "sync/atomic", "AddInt32", "AddInt64", "AddUint32", "AddUint64", "AddUintptr") {
				checkAtomicAddAssignment(pass, n.Lhs[i], call)
			}
		}
	})
	return nil, nil
}

// checkAtomicAddAssignment walks the atomic.Add* method calls checking
// for assigning the return value to the same variable being used in the
// operation
func checkAtomicAddAssignment(pass *analysis.Pass, left ast.Expr, call *ast.CallExpr) {
	if len(call.Args) != 2 {
		return
	}
	arg := call.Args[0]
	broken := false

	golangfmt := func(e ast.Expr) string { return analysisinternal.Format(pass.Fset, e) }

	if uarg, ok := arg.(*ast.UnaryExpr); ok && uarg.Op == token.AND {
		broken = golangfmt(left) == golangfmt(uarg.X)
	} else if star, ok := left.(*ast.StarExpr); ok {
		broken = golangfmt(star.X) == golangfmt(arg)
	}

	if broken {
		pass.ReportRangef(left, "direct assignment to atomic value")
	}
}
