// Copyright 2017 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"internal/testenv"
	"testing"
)

func TestInternalLinkerCgolangExec(t *testing.T) {
	testenv.MustHaveCGO(t)
	// N.B. the golang build explictly doesn't pass through
	// -asan/-msan/-race, so we don't care about those.
	testenv.MustInternalLink(t, testenv.SpecialBuildTypes{Cgolang: true})
	testGolangExec(t, true, false)
}

func TestExternalLinkerCgolangExec(t *testing.T) {
	testenv.MustHaveCGO(t)
	testGolangExec(t, true, true)
}

func TestCgolangLib(t *testing.T) {
	testenv.MustHaveCGO(t)
	testGolangLib(t, true)
}
