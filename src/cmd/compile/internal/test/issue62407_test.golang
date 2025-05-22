// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"reflect"
	"testing"
)

//golang:noinline
func foo() string { return "foofoofoofoofoofo" } // len == 17

//golang:noinline
func empty() string { return "" }

func TestConcatBytes(t *testing.T) {
	empty := empty()
	s := foo()
	tests := map[string]struct {
		golangt  []byte
		want []byte
	}{
		"two empty elements":                 {golangt: []byte(empty + empty), want: []byte{}},
		"two nonempty elements":              {golangt: []byte(s + s), want: append([]byte(foo()), foo()...)},
		"one empty and one nonempty element": {golangt: []byte(s + empty), want: []byte(foo())},
		"multiple empty elements":            {golangt: []byte(empty + empty + empty + empty + empty + empty), want: []byte{}},
		"multiple nonempty elements":         {golangt: []byte("1" + "2" + "3" + "4" + "5" + "6"), want: []byte("123456")},
	}

	for name, test := range tests {
		if !reflect.DeepEqual(test.golangt, test.want) {
			t.Errorf("[%s] golangt: %s, want: %s", name, test.golangt, test.want)
		}
	}
}

func TestConcatBytesAllocations(t *testing.T) {
	empty := empty()
	s := foo()
	tests := map[string]struct {
		f      func() []byte
		allocs float64
	}{
		"two empty elements":      {f: func() []byte { return []byte(empty + empty) }, allocs: 0},
		"multiple empty elements": {f: func() []byte { return []byte(empty + empty + empty + empty + empty + empty) }, allocs: 0},

		"two elements":                       {f: func() []byte { return []byte(s + s) }, allocs: 1},
		"three elements":                     {f: func() []byte { return []byte(s + s + s) }, allocs: 1},
		"four elements":                      {f: func() []byte { return []byte(s + s + s + s) }, allocs: 1},
		"five elements":                      {f: func() []byte { return []byte(s + s + s + s + s) }, allocs: 1},
		"one empty and one nonempty element": {f: func() []byte { return []byte(s + empty) }, allocs: 1},
		"two empty and two nonempty element": {f: func() []byte { return []byte(s + empty + s + empty) }, allocs: 1},
	}
	for name, test := range tests {
		allocs := testing.AllocsPerRun(100, func() { test.f() })
		if allocs != test.allocs {
			t.Errorf("concatbytes [%s]: %v allocs, want %v", name, allocs, test.allocs)
		}
	}
}
