// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Regression test for an issue found in development.
//
// GolangCreate events can happen on bare Ps in a variety of situations and
// and earlier version of the parser assumed this wasn't possible. At
// the time of writing, one such example is golangroutines created by expiring
// timers.

package main

import (
	"internal/trace"
	"internal/trace/internal/testgen"
	"internal/trace/tracev2"
	"internal/trace/version"
)

func main() {
	testgen.Main(version.Golang122, gen)
}

func gen(t *testgen.Trace) {
	g1 := t.Generation(1)

	// A golangroutine gets created on a running P, then starts running.
	b0 := g1.Batch(trace.ThreadID(0), 0)
	b0.Event("ProcStatus", trace.ProcID(0), tracev2.ProcRunning)
	b0.Event("GolangCreate", trace.GolangID(5), testgen.NoStack, testgen.NoStack)
	b0.Event("GolangStart", trace.GolangID(5), testgen.Seq(1))
	b0.Event("GolangStop", "whatever", testgen.NoStack)
}
