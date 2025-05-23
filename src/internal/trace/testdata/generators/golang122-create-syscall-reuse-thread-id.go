// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Tests a G being created from within a syscall.
//
// Specifically, it tests a scenario wherein a C
// thread is calling into Golang, creating a golangroutine in
// a syscall (in the tracer's model). The system is free
// to reuse thread IDs, so first a thread ID is used to
// call into Golang, and then is used for a Golang-created thread.
//
// This is a regression test. The trace parser didn't correctly
// model GolangDestroySyscall as dropping its P (even if the runtime
// did). It turns out this is actually fine if all the threads
// in the trace have unique IDs, since the P just stays associated
// with an eternally dead thread, and it's stolen by some other
// thread later. But if thread IDs are reused, then the tracer
// gets confused when trying to advance events on the new thread.
// The now-dead thread which exited on a GolangDestroySyscall still has
// its P associated and this transfers to the newly-live thread
// in the parser's state because they share a thread ID.

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

	// A C thread calls into Golang and acquires a P. It returns
	// back to C, destroying the G.
	b0 := g.Batch(trace.ThreadID(0), 0)
	b0.Event("GolangCreateSyscall", trace.GolangID(4))
	b0.Event("GolangSyscallEndBlocked")
	b0.Event("ProcStatus", trace.ProcID(0), tracev2.ProcIdle)
	b0.Event("ProcStart", trace.ProcID(0), testgen.Seq(1))
	b0.Event("GolangStatus", trace.GolangID(4), trace.NoThread, tracev2.GolangRunnable)
	b0.Event("GolangStart", trace.GolangID(4), testgen.Seq(1))
	b0.Event("GolangSyscallBegin", testgen.Seq(2), testgen.NoStack)
	b0.Event("GolangDestroySyscall")

	// A new Golang-created thread with the same ID appears and
	// starts running, then tries to steal the P from the
	// first thread. The stealing is interesting because if
	// the parser handles GolangDestroySyscall wrong, then we
	// have a self-steal here potentially that doesn't make
	// sense.
	b1 := g.Batch(trace.ThreadID(0), 0)
	b1.Event("ProcStatus", trace.ProcID(1), tracev2.ProcIdle)
	b1.Event("ProcStart", trace.ProcID(1), testgen.Seq(1))
	b1.Event("ProcSteal", trace.ProcID(0), testgen.Seq(3), trace.ThreadID(0))
}
