// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package noder

import (
	"reflect"
	"runtime"
	"testing"

	"cmd/compile/internal/syntax"
)

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestPragmaFields(t *testing.T) {
	var tests = []struct {
		in   string
		want []string
	}{
		{"", []string{}},
		{" \t ", []string{}},
		{`""""`, []string{`""`, `""`}},
		{"  a'b'c  ", []string{"a'b'c"}},
		{"1 2 3 4", []string{"1", "2", "3", "4"}},
		{"\n☺\t☹\n", []string{"☺", "☹"}},
		{`"1 2 "  3  " 4 5"`, []string{`"1 2 "`, `3`, `" 4 5"`}},
		{`"1""2 3""4"`, []string{`"1"`, `"2 3"`, `"4"`}},
		{`12"34"`, []string{`12`, `"34"`}},
		{`12"34 `, []string{`12`}},
	}

	for _, tt := range tests {
		golangt := pragmaFields(tt.in)
		if !eq(golangt, tt.want) {
			t.Errorf("pragmaFields(%q) = %v; want %v", tt.in, golangt, tt.want)
			continue
		}
	}
}

func TestPragcgolang(t *testing.T) {
	type testStruct struct {
		in   string
		want []string
	}

	var tests = []testStruct{
		{`golang:cgolang_export_dynamic local`, []string{`cgolang_export_dynamic`, `local`}},
		{`golang:cgolang_export_dynamic local remote`, []string{`cgolang_export_dynamic`, `local`, `remote`}},
		{`golang:cgolang_export_dynamic local' remote'`, []string{`cgolang_export_dynamic`, `local'`, `remote'`}},
		{`golang:cgolang_export_static local`, []string{`cgolang_export_static`, `local`}},
		{`golang:cgolang_export_static local remote`, []string{`cgolang_export_static`, `local`, `remote`}},
		{`golang:cgolang_export_static local' remote'`, []string{`cgolang_export_static`, `local'`, `remote'`}},
		{`golang:cgolang_import_dynamic local`, []string{`cgolang_import_dynamic`, `local`}},
		{`golang:cgolang_import_dynamic local remote`, []string{`cgolang_import_dynamic`, `local`, `remote`}},
		{`golang:cgolang_import_static local`, []string{`cgolang_import_static`, `local`}},
		{`golang:cgolang_import_static local'`, []string{`cgolang_import_static`, `local'`}},
		{`golang:cgolang_dynamic_linker "/path/"`, []string{`cgolang_dynamic_linker`, `/path/`}},
		{`golang:cgolang_dynamic_linker "/p ath/"`, []string{`cgolang_dynamic_linker`, `/p ath/`}},
		{`golang:cgolang_ldflag "arg"`, []string{`cgolang_ldflag`, `arg`}},
		{`golang:cgolang_ldflag "a rg"`, []string{`cgolang_ldflag`, `a rg`}},
	}

	if runtime.GOOS != "aix" {
		tests = append(tests, []testStruct{
			{`golang:cgolang_import_dynamic local remote "library"`, []string{`cgolang_import_dynamic`, `local`, `remote`, `library`}},
			{`golang:cgolang_import_dynamic local' remote' "lib rary"`, []string{`cgolang_import_dynamic`, `local'`, `remote'`, `lib rary`}},
		}...)
	} else {
		// cgolang_import_dynamic with a library is slightly different on AIX
		// as the library field must follow the pattern [libc.a/object.o].
		tests = append(tests, []testStruct{
			{`golang:cgolang_import_dynamic local remote "lib.a/obj.o"`, []string{`cgolang_import_dynamic`, `local`, `remote`, `lib.a/obj.o`}},
			// This test must fail.
			{`golang:cgolang_import_dynamic local' remote' "library"`, []string{`<unknown position>: usage: //golang:cgolang_import_dynamic local [remote ["lib.a/object.o"]]`}},
		}...)

	}

	var p noder
	var nopos syntax.Pos
	for _, tt := range tests {

		p.err = make(chan syntax.Error)
		golangtch := make(chan [][]string, 1)
		golang func() {
			p.pragcgolangbuf = nil
			p.pragcgolang(nopos, tt.in)
			if p.pragcgolangbuf != nil {
				golangtch <- p.pragcgolangbuf
			}
		}()

		select {
		case e := <-p.err:
			want := tt.want[0]
			if e.Error() != want {
				t.Errorf("pragcgolang(%q) = %q; want %q", tt.in, e, want)
				continue
			}
		case golangt := <-golangtch:
			want := [][]string{tt.want}
			if !reflect.DeepEqual(golangt, want) {
				t.Errorf("pragcgolang(%q) = %q; want %q", tt.in, golangt, want)
				continue
			}
		}

	}
}
