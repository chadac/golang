// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

/*
Trace is a tool for viewing trace files.

Trace files can be generated with:
  - runtime/trace.Start
  - net/http/pprof package
  - golang test -trace

Example usage:
Generate a trace file with 'golang test':

	golang test -trace trace.out pkg

View the trace in a web browser:

	golang tool trace trace.out

Generate a pprof-like profile from the trace:

	golang tool trace -pprof=TYPE trace.out > TYPE.pprof

Supported profile types are:
  - net: network blocking profile
  - sync: synchronization blocking profile
  - syscall: syscall blocking profile
  - sched: scheduler latency profile

Then, you can use the pprof tool to analyze the profile:

	golang tool pprof TYPE.pprof

Note that while the various profiles available when launching
'golang tool trace' work on every browser, the trace viewer itself
(the 'view trace' page) comes from the Chrome/Chromium project
and is only actively tested on that browser.
*/
package main
