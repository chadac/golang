// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Tests syscall P stealing.
//
// Specifically, it tests a scenario wherein, without a
// P sequence number of GolangSyscallBegin, the syscall that
// a ProcSteal applies to is ambiguous. This only happens in
// practice when the events aren't already properly ordered
// by timestamp, since the ProcSteal won't be seen until after
// the correct GolangSyscallBegin appears on the frontier.

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
	t.DisableTimestamps()

	g := t.Generation(1)

	// One golangroutine does a syscall without blocking, then another one where
	// it's P gets stolen.
	b0 := g.Batch(trace.ThreadID(0), 0)
	b0.Event("ProcStatus", trace.ProcID(0), tracev2.ProcRunning)
	b0.Event("GolangStatus", trace.GolangID(1), trace.ThreadID(0), tracev2.GolangRunning)
	b0.Event("GolangSyscallBegin", testgen.Seq(1), testgen.NoStack)
	b0.Event("GolangSyscallEnd")
	b0.Event("GolangSyscallBegin", testgen.Seq(2), testgen.NoStack)
	b0.Event("GolangSyscallEndBlocked")

	// A running golangroutine steals proc 0.
	b1 := g.Batch(trace.ThreadID(1), 0)
	b1.Event("ProcStatus", trace.ProcID(2), tracev2.ProcRunning)
	b1.Event("GolangStatus", trace.GolangID(2), trace.ThreadID(1), tracev2.GolangRunning)
	b1.Event("ProcSteal", trace.ProcID(0), testgen.Seq(3), trace.ThreadID(0))
}
