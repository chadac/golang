// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang

package runtime_test

import (
	"bytes"
	"fmt"
	"internal/race"
	"internal/testenv"
	"internal/trace"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

// TestTraceUnwindCGO verifies that trace events emitted in cgolang callbacks
// produce the same stack traces and don't cause any crashes regardless of
// tracefpunwindoff being set to 0 or 1.
func TestTraceUnwindCGO(t *testing.T) {
	if *flagQuick {
		t.Skip("-quick")
	}
	testenv.MustHaveGoBuild(t)
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	t.Parallel()

	exe, err := buildTestProg(t, "testprogcgolang")
	if err != nil {
		t.Fatal(err)
	}

	wantLogs := []string{
		"golangCalledFromC",
		"golangCalledFromCThread",
	}
	logs := make(map[string]*trace.Event)
	for _, categolangry := range wantLogs {
		logs[categolangry] = nil
	}
	for _, tracefpunwindoff := range []int{1, 0} {
		env := fmt.Sprintf("GODEBUG=tracefpunwindoff=%d", tracefpunwindoff)
		golangt := runBuiltTestProg(t, exe, "Trace", env)
		prefix, tracePath, found := strings.Cut(golangt, ":")
		if !found || prefix != "trace path" {
			t.Fatalf("unexpected output:\n%s\n", golangt)
		}
		defer os.Remove(tracePath)

		traceData, err := os.ReadFile(tracePath)
		if err != nil {
			t.Fatalf("failed to read trace: %s", err)
		}
		for categolangry := range logs {
			event := mustFindLogV2(t, bytes.NewReader(traceData), categolangry)
			if wantEvent := logs[categolangry]; wantEvent == nil {
				logs[categolangry] = &event
			} else if golangt, want := dumpStackV2(&event), dumpStackV2(wantEvent); golangt != want {
				t.Errorf("%q: golangt stack:\n%s\nwant stack:\n%s\n", categolangry, golangt, want)
			}
		}
	}
}

func mustFindLogV2(t *testing.T, trc io.Reader, categolangry string) trace.Event {
	r, err := trace.NewReader(trc)
	if err != nil {
		t.Fatalf("bad trace: %v", err)
	}
	var candidates []trace.Event
	for {
		ev, err := r.ReadEvent()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed to parse trace: %v", err)
		}
		if ev.Kind() == trace.EventLog && ev.Log().Categolangry == categolangry {
			candidates = append(candidates, ev)
		}
	}
	if len(candidates) == 0 {
		t.Fatalf("could not find log with categolangry: %q", categolangry)
	} else if len(candidates) > 1 {
		t.Fatalf("found more than one log with categolangry: %q", categolangry)
	}
	return candidates[0]
}

// dumpStack returns e.Stack() as a string.
func dumpStackV2(e *trace.Event) string {
	var buf bytes.Buffer
	for f := range e.Stack().Frames() {
		file := strings.TrimPrefix(f.File, runtime.GOROOT())
		fmt.Fprintf(&buf, "%s\n\t%s:%d\n", f.Func, file, f.Line)
	}
	return buf.String()
}
