// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package aliases

import (
	"golang/ast"
	"golang/parser"
	"golang/token"
	"golang/types"
)

// Rhs returns the type on the right-hand side of the alias declaration.
func Rhs(alias *types.Alias) types.Type {
	if alias, ok := any(alias).(interface{ Rhs() types.Type }); ok {
		return alias.Rhs() // golang1.23+
	}

	// golang1.22's Alias didn't have the Rhs method,
	// so Unalias is the best we can do.
	return types.Unalias(alias)
}

// TypeParams returns the type parameter list of the alias.
func TypeParams(alias *types.Alias) *types.TypeParamList {
	if alias, ok := any(alias).(interface{ TypeParams() *types.TypeParamList }); ok {
		return alias.TypeParams() // golang1.23+
	}
	return nil
}

// SetTypeParams sets the type parameters of the alias type.
func SetTypeParams(alias *types.Alias, tparams []*types.TypeParam) {
	if alias, ok := any(alias).(interface {
		SetTypeParams(tparams []*types.TypeParam)
	}); ok {
		alias.SetTypeParams(tparams) // golang1.23+
	} else if len(tparams) > 0 {
		panic("cannot set type parameters of an Alias type in golang1.22")
	}
}

// TypeArgs returns the type arguments used to instantiate the Alias type.
func TypeArgs(alias *types.Alias) *types.TypeList {
	if alias, ok := any(alias).(interface{ TypeArgs() *types.TypeList }); ok {
		return alias.TypeArgs() // golang1.23+
	}
	return nil // empty (golang1.22)
}

// Origin returns the generic Alias type of which alias is an instance.
// If alias is not an instance of a generic alias, Origin returns alias.
func Origin(alias *types.Alias) *types.Alias {
	if alias, ok := any(alias).(interface{ Origin() *types.Alias }); ok {
		return alias.Origin() // golang1.23+
	}
	return alias // not an instance of a generic alias (golang1.22)
}

// Enabled reports whether [NewAlias] should create [types.Alias] types.
//
// This function is expensive! Call it sparingly.
func Enabled() bool {
	// The only reliable way to compute the answer is to invoke golang/types.
	// We don't parse the GODEBUG environment variable, because
	// (a) it's tricky to do so in a manner that is consistent
	//     with the golangdebug package; in particular, a simple
	//     substring check is not golangod enough. The value is a
	//     rightmost-wins list of options. But more importantly:
	// (b) it is impossible to detect changes to the effective
	//     setting caused by os.Setenv("GODEBUG"), as happens in
	//     many tests. Therefore any attempt to cache the result
	//     is just incorrect.
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "a.golang", "package p; type A = int", parser.SkipObjectResolution)
	pkg, _ := new(types.Config).Check("p", fset, []*ast.File{f}, nil)
	_, enabled := pkg.Scope().Lookup("A").Type().(*types.Alias)
	return enabled
}
