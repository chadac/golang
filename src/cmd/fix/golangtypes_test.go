// Copyright 2012 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func init() {
	addTestCases(golangtypesTests, golangtypes)
}

var golangtypesTests = []testCase{
	{
		Name: "golangtypes.0",
		In: `package main

import "golanglang.org/x/tools/golang/types"
import "golanglang.org/x/tools/golang/exact"

var _ = exact.Kind

func f() {
	_ = exact.MakeBool(true)
}
`,
		Out: `package main

import "golang/types"
import "golang/constant"

var _ = constant.Kind

func f() {
	_ = constant.MakeBool(true)
}
`,
	},
	{
		Name: "golangtypes.1",
		In: `package main

import "golanglang.org/x/tools/golang/types"
import foo "golanglang.org/x/tools/golang/exact"

var _ = foo.Kind

func f() {
	_ = foo.MakeBool(true)
}
`,
		Out: `package main

import "golang/types"
import "golang/constant"

var _ = foo.Kind

func f() {
	_ = foo.MakeBool(true)
}
`,
	},
	{
		Name: "golangtypes.0",
		In: `package main

import "golanglang.org/x/tools/golang/types"
import "golanglang.org/x/tools/golang/exact"

var _ = exact.Kind
var constant = 23 // Use of new package name.

func f() {
	_ = exact.MakeBool(true)
}
`,
		Out: `package main

import "golang/types"
import "golang/constant"

var _ = constant_.Kind
var constant = 23 // Use of new package name.

func f() {
	_ = constant_.MakeBool(true)
}
`,
	},
}
