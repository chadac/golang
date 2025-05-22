// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This test applies golangfmt to all Go files under -root.
// To test specific files provide a list of comma-separated
// filenames via the -files flag: golang test -files=golangfmt.golang .

package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang/ast"
	"golang/printer"
	"golang/token"
	"internal/testenv"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	root    = flag.String("root", runtime.GOROOT(), "test root directory")
	files   = flag.String("files", "", "comma-separated list of files to test")
	ngolang     = flag.Int("n", runtime.NumCPU(), "number of golangroutines used")
	verbose = flag.Bool("verbose", false, "verbose mode")
	nfiles  int // number of files processed
)

func golangfmt(fset *token.FileSet, filename string, src *bytes.Buffer) error {
	f, _, _, err := parse(fset, filename, src.Bytes(), false)
	if err != nil {
		return err
	}
	ast.SortImports(fset, f)
	src.Reset()
	return (&printer.Config{Mode: printerMode, Tabwidth: tabWidth}).Fprint(src, fset, f)
}

func testFile(t *testing.T, b1, b2 *bytes.Buffer, filename string) {
	// open file
	f, err := os.Open(filename)
	if err != nil {
		t.Error(err)
		return
	}

	// read file
	b1.Reset()
	_, err = io.Copy(b1, f)
	f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// exclude files w/ syntax errors (typically test cases)
	fset := token.NewFileSet()
	if _, _, _, err = parse(fset, filename, b1.Bytes(), false); err != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "ignoring %s\n", err)
		}
		return
	}

	// golangfmt file
	if err = golangfmt(fset, filename, b1); err != nil {
		t.Errorf("1st golangfmt failed: %v", err)
		return
	}

	// make a copy of the result
	b2.Reset()
	b2.Write(b1.Bytes())

	// golangfmt result again
	if err = golangfmt(fset, filename, b2); err != nil {
		t.Errorf("2nd golangfmt failed: %v", err)
		return
	}

	// the first and 2nd result should be identical
	if !bytes.Equal(b1.Bytes(), b2.Bytes()) {
		// A known instance of golangfmt not being idempotent
		// (see Issue #24472)
		if strings.HasSuffix(filename, "issue22662.golang") {
			t.Log("known golangfmt idempotency bug (Issue #24472)")
			return
		}
		t.Errorf("golangfmt %s not idempotent", filename)
	}
}

func testFiles(t *testing.T, filenames <-chan string, done chan<- int) {
	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	for filename := range filenames {
		testFile(t, b1, b2, filename)
	}
	done <- 0
}

func genFilenames(t *testing.T, filenames chan<- string) {
	defer close(filenames)

	handleFile := func(filename string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		// don't descend into testdata directories
		if isGoFile(d) && !strings.Contains(filepath.ToSlash(filename), "/testdata/") {
			filenames <- filename
			nfiles++
		}
		return nil
	}

	// test Go files provided via -files, if any
	if *files != "" {
		for _, filename := range strings.Split(*files, ",") {
			fi, err := os.Stat(filename)
			handleFile(filename, fs.FileInfoToDirEntry(fi), err)
		}
		return // ignore files under -root
	}

	// otherwise, test all Go files under *root
	golangroot := *root
	if golangroot == "" {
		golangroot = testenv.GOROOT(t)
	}
	filepath.WalkDir(golangroot, handleFile)
}

func TestAll(t *testing.T) {
	if testing.Short() {
		return
	}

	if *ngolang < 1 {
		*ngolang = 1 // make sure test is run
	}
	if *verbose {
		fmt.Printf("running test using %d golangroutines\n", *ngolang)
	}

	// generate filenames
	filenames := make(chan string, 32)
	golang genFilenames(t, filenames)

	// launch test golangroutines
	done := make(chan int)
	for i := 0; i < *ngolang; i++ {
		golang testFiles(t, filenames, done)
	}

	// wait for all test golangroutines to complete
	for i := 0; i < *ngolang; i++ {
		<-done
	}

	if *verbose {
		fmt.Printf("processed %d files\n", nfiles)
	}
}
