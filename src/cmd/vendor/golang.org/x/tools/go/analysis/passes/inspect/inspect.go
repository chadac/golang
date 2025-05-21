// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package inspect defines an Analyzer that provides an AST inspector
// (golanglang.org/x/tools/golang/ast/inspector.Inspector) for the syntax trees
// of a package. It is only a building block for other analyzers.
//
// Example of use in another analysis:
//
//	import (
//		"golanglang.org/x/tools/golang/analysis"
//		"golanglang.org/x/tools/golang/analysis/passes/inspect"
//		"golanglang.org/x/tools/golang/ast/inspector"
//	)
//
//	var Analyzer = &analysis.Analyzer{
//		...
//		Requires:       []*analysis.Analyzer{inspect.Analyzer},
//	}
//
//	func run(pass *analysis.Pass) (interface{}, error) {
//		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
//		inspect.Preorder(nil, func(n ast.Node) {
//			...
//		})
//		return nil, nil
//	}
package inspect

import (
	"reflect"

	"golanglang.org/x/tools/golang/analysis"
	"golanglang.org/x/tools/golang/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:             "inspect",
	Doc:              "optimize AST traversal for later passes",
	URL:              "https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/inspect",
	Run:              run,
	RunDespiteErrors: true,
	ResultType:       reflect.TypeOf(new(inspector.Inspector)),
}

func run(pass *analysis.Pass) (any, error) {
	return inspector.New(pass.Files), nil
}
