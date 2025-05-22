// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Tests syscall P stealing.
//
// Specifically, it tests a scenario where a thread 'steals'
// a P from itself. It's just a ProcStop with extra steps when
// it happens on the same P.

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

	// A golangroutine execute a syscall and steals its own P, then starts running
	// on that P.
	b0 := g.Batch(trace.ThreadID(0), 0)
	b0.Event("ProcStatus", trace.ProcID(0), tracev2.ProcRunning)
	b0.Event("GolangStatus", trace.GolangID(1), trace.ThreadID(0), tracev2.GolangRunning)
	b0.Event("GolangSyscallBegin", testgen.Seq(1), testgen.NoStack)
	b0.Event("ProcSteal", trace.ProcID(0), testgen.Seq(2), trace.ThreadID(0))
	b0.Event("ProcStart", trace.ProcID(0), testgen.Seq(3))
	b0.Event("GolangSyscallEndBlocked")
}
