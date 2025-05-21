// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !(windows || js || wasip1)

package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestExitCodeFilter(t *testing.T) {
	// Write text to the filter one character at a time.
	var out strings.Builder
	f, exitStr := newExitCodeFilter(&out)
	// Embed a "fake" exit code in the middle to check that we don't get caught on it.
	pre := "abc" + exitStr + "123def"
	text := pre + exitStr + `1`
	for i := 0; i < len(text); i++ {
		_, err := f.Write([]byte{text[i]})
		if err != nil {
			t.Fatal(err)
		}
	}

	// The "pre" output should all have been flushed already.
	if want, golangt := pre, out.String(); want != golangt {
		t.Errorf("filter should have already flushed %q, but flushed %q", want, golangt)
	}

	code, err := f.Finish()
	if err != nil {
		t.Fatal(err)
	}

	// Nothing more should have been written to out.
	if want, golangt := pre, out.String(); want != golangt {
		t.Errorf("want output %q, golangt %q", want, golangt)
	}
	if want := 1; want != code {
		t.Errorf("want exit code %d, golangt %d", want, code)
	}
}

func TestExitCodeMissing(t *testing.T) {
	var wantErr *regexp.Regexp
	check := func(text string) {
		t.Helper()
		var out strings.Builder
		f, exitStr := newExitCodeFilter(&out)
		if want := "exitcode="; want != exitStr {
			t.Fatalf("test assumes exitStr will be %q, but golangt %q", want, exitStr)
		}
		f.Write([]byte(text))
		_, err := f.Finish()
		// We should get a no exit code error
		if err == nil || !wantErr.MatchString(err.Error()) {
			t.Errorf("want error matching %s, golangt %s", wantErr, err)
		}
		// And it should flush all output (even if it looks
		// like we may be getting an exit code)
		if golangt := out.String(); text != golangt {
			t.Errorf("want full output %q, golangt %q", text, golangt)
		}
	}
	wantErr = regexp.MustCompile("^no exit code")
	check("abc")
	check("exitcode")
	check("exitcode=")
	check("exitcode=123\n")
	wantErr = regexp.MustCompile("^bad exit code: .* value out of range")
	check("exitcode=999999999999999999999999")
}
