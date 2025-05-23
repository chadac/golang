// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Tests syscall P stealing from a golangroutine and thread
// that have been in a syscall the entire generation.

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
	g := t.Generation(1)

	// Steal proc from a golangroutine that's been blocked
	// in a syscall the entire generation.
	b0 := g.Batch(trace.ThreadID(0), 0)
	b0.Event("ProcStatus", trace.ProcID(0), tracev2.ProcSyscallAbandoned)
	b0.Event("ProcSteal", trace.ProcID(0), testgen.Seq(1), trace.ThreadID(1))

	// Status event for a golangroutine blocked in a syscall for the entire generation.
	bz := g.Batch(trace.NoThread, 0)
	bz.Event("GolangStatus", trace.GolangID(1), trace.ThreadID(1), tracev2.GolangSyscall)
}
