// Copyright 2021 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build compiler_bootstrap

package main

import (
	"golang/ast"
	"golang/token"
)

func (f *File) walkUnexpected(x interface{}, context astContext, visit func(*File, interface{}, astContext)) {
	error_(token.NoPos, "unexpected type %T in walk", x)
	panic("unexpected type")
}

func funcTypeTypeParams(n *ast.FuncType) *ast.FieldList {
	return nil
}

func typeSpecTypeParams(n *ast.TypeSpec) *ast.FieldList {
	return nil
}
