// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"internal/trace"
)

var _ generator = &golangroutineGenerator{}

type golangroutineGenerator struct {
	globalRangeGenerator
	globalMetricGenerator
	stackSampleGenerator[trace.GolangID]
	logEventGenerator[trace.GolangID]

	gStates map[trace.GolangID]*gState[trace.GolangID]
	focus   trace.GolangID
	filter  map[trace.GolangID]struct{}
}

func newGolangroutineGenerator(ctx *traceContext, focus trace.GolangID, filter map[trace.GolangID]struct{}) *golangroutineGenerator {
	gg := new(golangroutineGenerator)
	rg := func(ev *trace.Event) trace.GolangID {
		return ev.Golangroutine()
	}
	gg.stackSampleGenerator.getResource = rg
	gg.logEventGenerator.getResource = rg
	gg.gStates = make(map[trace.GolangID]*gState[trace.GolangID])
	gg.focus = focus
	gg.filter = filter

	// Enable a filter on the emitter.
	if filter != nil {
		ctx.SetResourceFilter(func(resource uint64) bool {
			_, ok := filter[trace.GolangID(resource)]
			return ok
		})
	}
	return gg
}

func (g *golangroutineGenerator) Sync() {
	g.globalRangeGenerator.Sync()
}

func (g *golangroutineGenerator) GolangroutineLabel(ctx *traceContext, ev *trace.Event) {
	l := ev.Label()
	g.gStates[l.Resource.Golangroutine()].setLabel(l.Label)
}

func (g *golangroutineGenerator) GolangroutineRange(ctx *traceContext, ev *trace.Event) {
	r := ev.Range()
	switch ev.Kind() {
	case trace.EventRangeBegin:
		g.gStates[r.Scope.Golangroutine()].rangeBegin(ev.Time(), r.Name, ev.Stack())
	case trace.EventRangeActive:
		g.gStates[r.Scope.Golangroutine()].rangeActive(r.Name)
	case trace.EventRangeEnd:
		gs := g.gStates[r.Scope.Golangroutine()]
		gs.rangeEnd(ev.Time(), r.Name, ev.Stack(), ctx)
	}
}

func (g *golangroutineGenerator) GolangroutineTransition(ctx *traceContext, ev *trace.Event) {
	st := ev.StateTransition()
	golangID := st.Resource.Golangroutine()

	// If we haven't seen this golangroutine before, create a new
	// gState for it.
	gs, ok := g.gStates[golangID]
	if !ok {
		gs = newGState[trace.GolangID](golangID)
		g.gStates[golangID] = gs
	}

	// Try to augment the name of the golangroutine.
	gs.augmentName(st.Stack)

	// Handle the golangroutine state transition.
	from, to := st.Golangroutine()
	if from == to {
		// Filter out no-op events.
		return
	}
	if from.Executing() && !to.Executing() {
		if to == trace.GolangWaiting {
			// Golangroutine started blocking.
			gs.block(ev.Time(), ev.Stack(), st.Reason, ctx)
		} else {
			gs.stop(ev.Time(), ev.Stack(), ctx)
		}
	}
	if !from.Executing() && to.Executing() {
		start := ev.Time()
		if from == trace.GolangUndetermined {
			// Back-date the event to the start of the trace.
			start = ctx.startTime
		}
		gs.start(start, golangID, ctx)
	}

	if from == trace.GolangWaiting {
		// Golangroutine unblocked.
		gs.unblock(ev.Time(), ev.Stack(), ev.Golangroutine(), ctx)
	}
	if from == trace.GolangNotExist && to == trace.GolangRunnable {
		// Golangroutine was created.
		gs.created(ev.Time(), ev.Golangroutine(), ev.Stack())
	}
	if from == trace.GolangSyscall && to != trace.GolangRunning {
		// Exiting blocked syscall.
		gs.syscallEnd(ev.Time(), true, ctx)
		gs.blockedSyscallEnd(ev.Time(), ev.Stack(), ctx)
	} else if from == trace.GolangSyscall {
		// Check if we're exiting a syscall in a non-blocking way.
		gs.syscallEnd(ev.Time(), false, ctx)
	}

	// Handle syscalls.
	if to == trace.GolangSyscall {
		start := ev.Time()
		if from == trace.GolangUndetermined {
			// Back-date the event to the start of the trace.
			start = ctx.startTime
		}
		// Write down that we've entered a syscall. Note: we might have no G or P here
		// if we're in a cgolang callback or this is a transition from GolangUndetermined
		// (i.e. the G has been blocked in a syscall).
		gs.syscallBegin(start, golangID, ev.Stack())
	}

	// Note down the golangroutine transition.
	_, inMarkAssist := gs.activeRanges["GC mark assist"]
	ctx.GolangroutineTransition(ctx.elapsed(ev.Time()), viewerGState(from, inMarkAssist), viewerGState(to, inMarkAssist))
}

func (g *golangroutineGenerator) ProcRange(ctx *traceContext, ev *trace.Event) {
	// TODO(mknyszek): Extend procRangeGenerator to support rendering proc ranges
	// that overlap with a golangroutine's execution.
}

func (g *golangroutineGenerator) ProcTransition(ctx *traceContext, ev *trace.Event) {
	// Not needed. All relevant information for golangroutines can be derived from golangroutine transitions.
}

func (g *golangroutineGenerator) Finish(ctx *traceContext) {
	ctx.SetResourceType("G")

	// Finish off global ranges.
	g.globalRangeGenerator.Finish(ctx)

	// Finish off all the golangroutine slices.
	for id, gs := range g.gStates {
		gs.finish(ctx)

		// Tell the emitter about the golangroutines we want to render.
		ctx.Resource(uint64(id), gs.name())
	}

	// Set the golangroutine to focus on.
	if g.focus != trace.NoGolangroutine {
		ctx.Focus(uint64(g.focus))
	}
}
