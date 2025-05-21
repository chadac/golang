// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

func init() {
	addTestCases(buildtagTests, buildtag)
}

var buildtagTests = []testCase{
	{
		Name:    "buildtag.oldGo",
		Version: "golang1.10",
		In: `//golang:build yes
// +build yes

package main
`,
	},
	{
		Name:    "buildtag.new",
		Version: "golang1.99",
		In: `//golang:build yes
// +build yes

package main
`,
		Out: `//golang:build yes

package main
`,
	},
}
