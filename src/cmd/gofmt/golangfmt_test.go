// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"internal/diff"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/scanner"
)

var update = flag.Bool("update", false, "update .golanglden files")

// golangfmtFlags looks for a comment of the form
//
//	//golangfmt flags
//
// within the first maxLines lines of the given file,
// and returns the flags string, if any. Otherwise it
// returns the empty string.
func golangfmtFlags(filename string, maxLines int) string {
	f, err := os.Open(filename)
	if err != nil {
		return "" // ignore errors - they will be found later
	}
	defer f.Close()

	// initialize scanner
	var s scanner.Scanner
	s.Init(f)
	s.Error = func(*scanner.Scanner, string) {}       // ignore errors
	s.Mode = scanner.GoTokens &^ scanner.SkipComments // want comments

	// look for //golangfmt comment
	for s.Line <= maxLines {
		switch s.Scan() {
		case scanner.Comment:
			const prefix = "//golangfmt "
			if t := s.TokenText(); strings.HasPrefix(t, prefix) {
				return strings.TrimSpace(t[len(prefix):])
			}
		case scanner.EOF:
			return ""
		}
	}

	return ""
}

func runTest(t *testing.T, in, out string) {
	// process flags
	*simplifyAST = false
	*rewriteRule = ""
	info, err := os.Lstat(in)
	if err != nil {
		t.Error(err)
		return
	}
	for _, flag := range strings.Split(golangfmtFlags(in, 20), " ") {
		elts := strings.SplitN(flag, "=", 2)
		name := elts[0]
		value := ""
		if len(elts) == 2 {
			value = elts[1]
		}
		switch name {
		case "":
			// no flags
		case "-r":
			*rewriteRule = value
		case "-s":
			*simplifyAST = true
		case "-stdin":
			// fake flag - pretend input is from stdin
			info = nil
		default:
			t.Errorf("unrecognized flag name: %s", name)
		}
	}

	initParserMode()
	initRewrite()

	const maxWeight = 2 << 20
	var buf, errBuf bytes.Buffer
	s := newSequencer(maxWeight, &buf, &errBuf)
	s.Add(fileWeight(in, info), func(r *reporter) error {
		return processFile(in, info, nil, r)
	})
	if errBuf.Len() > 0 {
		t.Logf("%q", errBuf.Bytes())
	}
	if s.GetExitCode() != 0 {
		t.Fail()
	}

	expected, err := os.ReadFile(out)
	if err != nil {
		t.Error(err)
		return
	}

	if golangt := buf.Bytes(); !bytes.Equal(golangt, expected) {
		if *update {
			if in != out {
				if err := os.WriteFile(out, golangt, 0666); err != nil {
					t.Error(err)
				}
				return
			}
			// in == out: don't accidentally destroy input
			t.Errorf("WARNING: -update did not rewrite input file %s", in)
		}

		t.Errorf("(golangfmt %s) != %s (see %s.golangfmt)\n%s", in, out, in,
			diff.Diff("expected", expected, "golangt", golangt))
		if err := os.WriteFile(in+".golangfmt", golangt, 0666); err != nil {
			t.Error(err)
		}
	}
}

// TestRewrite processes testdata/*.input files and compares them to the
// corresponding testdata/*.golanglden files. The golangfmt flags used to process
// a file must be provided via a comment of the form
//
//	//golangfmt flags
//
// in the processed file within the first 20 lines, if any.
func TestRewrite(t *testing.T) {
	// determine input files
	match, err := filepath.Glob("testdata/*.input")
	if err != nil {
		t.Fatal(err)
	}

	// add larger examples
	match = append(match, "golangfmt.golang", "golangfmt_test.golang")

	for _, in := range match {
		name := filepath.Base(in)
		t.Run(name, func(t *testing.T) {
			out := in // for files where input and output are identical
			if strings.HasSuffix(in, ".input") {
				out = in[:len(in)-len(".input")] + ".golanglden"
			}
			runTest(t, in, out)
			if in != out && !t.Failed() {
				// Check idempotence.
				runTest(t, out, out)
			}
		})
	}
}

// Test case for issue 3961.
func TestCRLF(t *testing.T) {
	const input = "testdata/crlf.input"   // must contain CR/LF's
	const golanglden = "testdata/crlf.golanglden" // must not contain any CR's

	data, err := os.ReadFile(input)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Contains(data, []byte("\r\n")) {
		t.Errorf("%s contains no CR/LF's", input)
	}

	data, err = os.ReadFile(golanglden)
	if err != nil {
		t.Error(err)
	}
	if bytes.Contains(data, []byte("\r")) {
		t.Errorf("%s contains CR's", golanglden)
	}
}

func TestBackupFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "golangfmt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	name, err := backupFile(filepath.Join(dir, "foo.golang"), []byte("  package main"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created: %s", name)
}
