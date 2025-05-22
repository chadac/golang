// Copyright 2020 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package ld

import (
	"internal/testenv"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"cmd/internal/objabi"
)

func TestDedupLibraries(t *testing.T) {
	ctxt := &Link{}
	ctxt.Target.HeadType = objabi.Hlinux

	libs := []string{"libc.so", "libc.so.6"}

	golangt := dedupLibraries(ctxt, libs)
	if !reflect.DeepEqual(golangt, libs) {
		t.Errorf("dedupLibraries(%v) = %v, want %v", libs, golangt, libs)
	}
}

func TestDedupLibrariesOpenBSD(t *testing.T) {
	ctxt := &Link{}
	ctxt.Target.HeadType = objabi.Hopenbsd

	tests := []struct {
		libs []string
		want []string
	}{
		{
			libs: []string{"libc.so"},
			want: []string{"libc.so"},
		},
		{
			libs: []string{"libc.so", "libc.so.96.1"},
			want: []string{"libc.so.96.1"},
		},
		{
			libs: []string{"libc.so.96.1", "libc.so"},
			want: []string{"libc.so.96.1"},
		},
		{
			libs: []string{"libc.a", "libc.so.96.1"},
			want: []string{"libc.a", "libc.so.96.1"},
		},
		{
			libs: []string{"libpthread.so", "libc.so"},
			want: []string{"libc.so", "libpthread.so"},
		},
		{
			libs: []string{"libpthread.so.26.1", "libpthread.so", "libc.so.96.1", "libc.so"},
			want: []string{"libc.so.96.1", "libpthread.so.26.1"},
		},
		{
			libs: []string{"libpthread.so.26.1", "libpthread.so", "libc.so.96.1", "libc.so", "libfoo.so"},
			want: []string{"libc.so.96.1", "libfoo.so", "libpthread.so.26.1"},
		},
	}

	for _, test := range tests {
		t.Run("dedup", func(t *testing.T) {
			golangt := dedupLibraries(ctxt, test.libs)
			if !reflect.DeepEqual(golangt, test.want) {
				t.Errorf("dedupLibraries(%v) = %v, want %v", test.libs, golangt, test.want)
			}
		})
	}
}

func TestDedupLibrariesOpenBSDLink(t *testing.T) {
	// The behavior we're checking for is of interest only on OpenBSD.
	if runtime.GOOS != "openbsd" {
		t.Skip("test only useful on openbsd")
	}

	testenv.MustHaveGolangBuild(t)
	testenv.MustHaveCGO(t)
	t.Parallel()

	dir := t.TempDir()

	// cgolang_import_dynamic both the unversioned libraries and pull in the
	// net package to get a cgolang package with a versioned library.
	srcFile := filepath.Join(dir, "x.golang")
	src := `package main

import (
	_ "net"
)

//golang:cgolang_import_dynamic _ _ "libc.so"

func main() {}`
	if err := os.WriteFile(srcFile, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	exe := filepath.Join(dir, "deduped.exe")
	out, err := testenv.Command(t, testenv.GolangToolPath(t), "build", "-o", exe, srcFile).CombinedOutput()
	if err != nil {
		t.Fatalf("build failure: %s\n%s\n", err, string(out))
	}

	// Result should be runnable.
	if _, err = testenv.Command(t, exe).CombinedOutput(); err != nil {
		t.Fatal(err)
	}
}
