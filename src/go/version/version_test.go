// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) { test2(t, compareTests, "Compare", Compare) }

var compareTests = []testCase2[string, string, int]{
	{"", "", 0},
	{"x", "x", 0},
	{"", "x", 0},
	{"1", "1.1", 0},
	{"golang1", "golang1.1", -1},
	{"golang1.5", "golang1.6", -1},
	{"golang1.5", "golang1.10", -1},
	{"golang1.6", "golang1.6.1", -1},
	{"golang1.19", "golang1.19.0", 0},
	{"golang1.19rc1", "golang1.19", -1},
	{"golang1.20", "golang1.20.0", 0},
	{"golang1.20", "golang1.20.0-bigcorp", 0},
	{"golang1.20rc1", "golang1.20", -1},
	{"golang1.21", "golang1.21.0", -1},
	{"golang1.21", "golang1.21.0-bigcorp", -1},
	{"golang1.21", "golang1.21rc1", -1},
	{"golang1.21rc1", "golang1.21.0", -1},
	{"golang1.6", "golang1.19", -1},
	{"golang1.19", "golang1.19.1", -1},
	{"golang1.19rc1", "golang1.19", -1},
	{"golang1.19rc1", "golang1.19", -1},
	{"golang1.19rc1", "golang1.19.1", -1},
	{"golang1.19rc1", "golang1.19rc2", -1},
	{"golang1.19.0", "golang1.19.1", -1},
	{"golang1.19rc1", "golang1.19.0", -1},
	{"golang1.19alpha3", "golang1.19beta2", -1},
	{"golang1.19beta2", "golang1.19rc1", -1},
	{"golang1.1", "golang1.99999999999999998", -1},
	{"golang1.99999999999999998", "golang1.99999999999999999", -1},
}

func TestLang(t *testing.T) { test1(t, langTests, "Lang", Lang) }

var langTests = []testCase1[string, string]{
	{"bad", ""},
	{"golang1.2rc3", "golang1.2"},
	{"golang1.2.3", "golang1.2"},
	{"golang1.2", "golang1.2"},
	{"golang1", "golang1"},
	{"golang222", "golang222.0"},
	{"golang1.999testmod", "golang1.999"},
}

func TestIsValid(t *testing.T) { test1(t, isValidTests, "IsValid", IsValid) }

var isValidTests = []testCase1[string, bool]{
	{"", false},
	{"1.2.3", false},
	{"golang1.2rc3", true},
	{"golang1.2.3", true},
	{"golang1.999testmod", true},
	{"golang1.600+auto", false},
	{"golang1.22", true},
	{"golang1.21.0", true},
	{"golang1.21rc2", true},
	{"golang1.21", true},
	{"golang1.20.0", true},
	{"golang1.20", true},
	{"golang1.19", true},
	{"golang1.3", true},
	{"golang1.2", true},
	{"golang1", true},
}

type testCase1[In, Out any] struct {
	in  In
	out Out
}

type testCase2[In1, In2, Out any] struct {
	in1 In1
	in2 In2
	out Out
}

func test1[In, Out any](t *testing.T, tests []testCase1[In, Out], name string, f func(In) Out) {
	t.Helper()
	for _, tt := range tests {
		if out := f(tt.in); !reflect.DeepEqual(out, tt.out) {
			t.Errorf("%s(%v) = %v, want %v", name, tt.in, out, tt.out)
		}
	}
}

func test2[In1, In2, Out any](t *testing.T, tests []testCase2[In1, In2, Out], name string, f func(In1, In2) Out) {
	t.Helper()
	for _, tt := range tests {
		if out := f(tt.in1, tt.in2); !reflect.DeepEqual(out, tt.out) {
			t.Errorf("%s(%+v, %+v) = %+v, want %+v", name, tt.in1, tt.in2, out, tt.out)
		}
	}
}
