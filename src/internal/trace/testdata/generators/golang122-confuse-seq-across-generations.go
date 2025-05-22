// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Regression test for an issue found in development.
//
// The core of the issue is that if generation counters
// aren't considered as part of sequence numbers, then
// it's possible to accidentally advance without a
// GolangStatus event.
//
// The situation is one in which it just so happens that
// an event on the frontier for a following generation
// has a sequence number exactly one higher than the last
// sequence number for e.g. a golangroutine in the previous
// generation. The parser should wait to find a GolangStatus
// event before advancing into the next generation at all.
// It turns out this situation is pretty rare; the GolangStatus
// event almost always shows up first in practice. But it
// can and did happen.

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

	// A running golangroutine blocks.
	b10 := g1.Batch(trace.ThreadID(0), 0)
	b10.Event("ProcStatus", trace.ProcID(0), tracev2.ProcRunning)
	b10.Event("GolangStatus", trace.GolangID(1), trace.ThreadID(0), tracev2.GolangRunning)
	b10.Event("GolangStop", "whatever", testgen.NoStack)

	// The running golangroutine gets unblocked.
	b11 := g1.Batch(trace.ThreadID(1), 0)
	b11.Event("ProcStatus", trace.ProcID(1), tracev2.ProcRunning)
	b11.Event("GolangStart", trace.GolangID(1), testgen.Seq(1))
	b11.Event("GolangStop", "whatever", testgen.NoStack)

	g2 := t.Generation(2)

	// Start running the golangroutine, but later.
	b21 := g2.Batch(trace.ThreadID(1), 3)
	b21.Event("ProcStatus", trace.ProcID(1), tracev2.ProcRunning)
	b21.Event("GolangStart", trace.GolangID(1), testgen.Seq(2))

	// The golangroutine starts running, then stops, then starts again.
	b20 := g2.Batch(trace.ThreadID(0), 5)
	b20.Event("ProcStatus", trace.ProcID(0), tracev2.ProcRunning)
	b20.Event("GolangStatus", trace.GolangID(1), trace.ThreadID(0), tracev2.GolangRunnable)
	b20.Event("GolangStart", trace.GolangID(1), testgen.Seq(1))
	b20.Event("GolangStop", "whatever", testgen.NoStack)
}
