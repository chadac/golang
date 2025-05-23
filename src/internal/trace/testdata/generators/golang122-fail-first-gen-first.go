// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Regression test for #55160.
//
// The issue is that the parser reads ahead to the first batch of the
// next generation to find generation boundaries, but if it finds an
// error, it needs to delay handling that error until later. Previously
// it would handle that error immediately and a totally valid generation
// would be skipped for parsing and rejected because of an error in a
// batch in the following generation.
//
// This test captures this behavior by making both the first generation
// and second generation bad. It requires that the issue in the first
// generation, which is caught when actually ordering events, be reported
// instead of the second one.

package main

import (
	"internal/trace/internal/testgen"
	"internal/trace/tracev2"
	"internal/trace/version"
)

func main() {
	testgen.Main(version.Golang122, gen)
}

func gen(t *testgen.Trace) {
	// A running golangroutine emits a task begin.
	t.RawEvent(tracev2.EvEventBatch, nil, 1 /*gen*/, 0 /*thread ID*/, 0 /*timestamp*/, 5 /*batch length*/)
	t.RawEvent(tracev2.EvFrequency, nil, 15625000)

	// A running golangroutine emits a task begin.
	t.RawEvent(tracev2.EvEventBatch, nil, 1 /*gen*/, 0 /*thread ID*/, 0 /*timestamp*/, 5 /*batch length*/)
	t.RawEvent(tracev2.EvGolangCreate, nil, 0 /*timestamp delta*/, 1 /*golang ID*/, 0, 0)

	// Write an invalid batch event for the next generation.
	t.RawEvent(tracev2.EvEventBatch, nil, 2 /*gen*/, 0 /*thread ID*/, 0 /*timestamp*/, 50 /*batch length (invalid)*/)

	// We should fail at the first issue, not the second one.
	t.ExpectFailure("expected a proc but didn't have one")
}
