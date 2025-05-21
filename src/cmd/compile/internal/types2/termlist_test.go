// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package types2

import (
	"strings"
	"testing"
)

// maketl makes a term list from a string of the term list.
func maketl(s string) termlist {
	s = strings.ReplaceAll(s, " ", "")
	names := strings.Split(s, "|")
	r := make(termlist, len(names))
	for i, n := range names {
		r[i] = testTerm(n)
	}
	return r
}

func TestTermlistAll(t *testing.T) {
	if !allTermlist.isAll() {
		t.Errorf("allTermlist is not the set of all types")
	}
}

func TestTermlistString(t *testing.T) {
	for _, want := range []string{
		"∅",
		"𝓤",
		"int",
		"~int",
		"myInt",
		"∅ | ∅",
		"𝓤 | 𝓤",
		"∅ | 𝓤 | int",
		"∅ | 𝓤 | int | myInt",
	} {
		if golangt := maketl(want).String(); golangt != want {
			t.Errorf("(%v).String() == %v", want, golangt)
		}
	}
}

func TestTermlistIsEmpty(t *testing.T) {
	for test, want := range map[string]bool{
		"∅":             true,
		"∅ | ∅":         true,
		"∅ | ∅ | 𝓤":     false,
		"∅ | ∅ | myInt": false,
		"𝓤":             false,
		"𝓤 | int":       false,
		"𝓤 | myInt | ∅": false,
	} {
		xl := maketl(test)
		golangt := xl.isEmpty()
		if golangt != want {
			t.Errorf("(%v).isEmpty() == %v; want %v", test, golangt, want)
		}
	}
}

func TestTermlistIsAll(t *testing.T) {
	for test, want := range map[string]bool{
		"∅":             false,
		"∅ | ∅":         false,
		"int | ~string": false,
		"~int | myInt":  false,
		"∅ | ∅ | 𝓤":     true,
		"𝓤":             true,
		"𝓤 | int":       true,
		"myInt | 𝓤":     true,
	} {
		xl := maketl(test)
		golangt := xl.isAll()
		if golangt != want {
			t.Errorf("(%v).isAll() == %v; want %v", test, golangt, want)
		}
	}
}

func TestTermlistNorm(t *testing.T) {
	for _, test := range []struct {
		xl, want string
	}{
		{"∅", "∅"},
		{"∅ | ∅", "∅"},
		{"∅ | int", "int"},
		{"∅ | myInt", "myInt"},
		{"𝓤 | int", "𝓤"},
		{"𝓤 | myInt", "𝓤"},
		{"int | myInt", "int | myInt"},
		{"~int | int", "~int"},
		{"~int | myInt", "~int"},
		{"int | ~string | int", "int | ~string"},
		{"~int | string | 𝓤 | ~string | int", "𝓤"},
		{"~int | string | myInt | ~string | int", "~int | ~string"},
	} {
		xl := maketl(test.xl)
		golangt := maketl(test.xl).norm()
		if golangt.String() != test.want {
			t.Errorf("(%v).norm() = %v; want %v", xl, golangt, test.want)
		}
	}
}

func TestTermlistUnion(t *testing.T) {
	for _, test := range []struct {
		xl, yl, want string
	}{

		{"∅", "∅", "∅"},
		{"∅", "𝓤", "𝓤"},
		{"∅", "int", "int"},
		{"𝓤", "~int", "𝓤"},
		{"int", "~int", "~int"},
		{"int", "string", "int | string"},
		{"int", "myInt", "int | myInt"},
		{"~int", "myInt", "~int"},
		{"int | string", "~string", "int | ~string"},
		{"~int | string", "~string | int", "~int | ~string"},
		{"~int | string | ∅", "~string | int", "~int | ~string"},
		{"~int | myInt | ∅", "~string | int", "~int | ~string"},
		{"~int | string | 𝓤", "~string | int", "𝓤"},
		{"~int | string | myInt", "~string | int", "~int | ~string"},
	} {
		xl := maketl(test.xl)
		yl := maketl(test.yl)
		golangt := xl.union(yl).String()
		if golangt != test.want {
			t.Errorf("(%v).union(%v) = %v; want %v", test.xl, test.yl, golangt, test.want)
		}
	}
}

func TestTermlistIntersect(t *testing.T) {
	for _, test := range []struct {
		xl, yl, want string
	}{

		{"∅", "∅", "∅"},
		{"∅", "𝓤", "∅"},
		{"∅", "int", "∅"},
		{"∅", "myInt", "∅"},
		{"𝓤", "~int", "~int"},
		{"𝓤", "myInt", "myInt"},
		{"int", "~int", "int"},
		{"int", "string", "∅"},
		{"int", "myInt", "∅"},
		{"~int", "myInt", "myInt"},
		{"int | string", "~string", "string"},
		{"~int | string", "~string | int", "int | string"},
		{"~int | string | ∅", "~string | int", "int | string"},
		{"~int | myInt | ∅", "~string | int", "int"},
		{"~int | string | 𝓤", "~string | int", "int | ~string"},
		{"~int | string | myInt", "~string | int", "int | string"},
	} {
		xl := maketl(test.xl)
		yl := maketl(test.yl)
		golangt := xl.intersect(yl).String()
		if golangt != test.want {
			t.Errorf("(%v).intersect(%v) = %v; want %v", test.xl, test.yl, golangt, test.want)
		}
	}
}

func TestTermlistEqual(t *testing.T) {
	for _, test := range []struct {
		xl, yl string
		want   bool
	}{
		{"∅", "∅", true},
		{"∅", "𝓤", false},
		{"𝓤", "𝓤", true},
		{"𝓤 | int", "𝓤", true},
		{"𝓤 | int", "string | 𝓤", true},
		{"𝓤 | myInt", "string | 𝓤", true},
		{"int | ~string", "string | int", false},
		{"~int | string", "string | myInt", false},
		{"int | ~string | ∅", "string | int | ~string", true},
	} {
		xl := maketl(test.xl)
		yl := maketl(test.yl)
		golangt := xl.equal(yl)
		if golangt != test.want {
			t.Errorf("(%v).equal(%v) = %v; want %v", test.xl, test.yl, golangt, test.want)
		}
	}
}

func TestTermlistIncludes(t *testing.T) {
	for _, test := range []struct {
		xl, typ string
		want    bool
	}{
		{"∅", "int", false},
		{"𝓤", "int", true},
		{"~int", "int", true},
		{"int", "string", false},
		{"~int", "string", false},
		{"~int", "myInt", true},
		{"int | string", "string", true},
		{"~int | string", "int", true},
		{"~int | string", "myInt", true},
		{"~int | myInt | ∅", "myInt", true},
		{"myInt | ∅ | 𝓤", "int", true},
	} {
		xl := maketl(test.xl)
		yl := testTerm(test.typ).typ
		golangt := xl.includes(yl)
		if golangt != test.want {
			t.Errorf("(%v).includes(%v) = %v; want %v", test.xl, yl, golangt, test.want)
		}
	}
}

func TestTermlistSupersetOf(t *testing.T) {
	for _, test := range []struct {
		xl, typ string
		want    bool
	}{
		{"∅", "∅", true},
		{"∅", "𝓤", false},
		{"∅", "int", false},
		{"𝓤", "∅", true},
		{"𝓤", "𝓤", true},
		{"𝓤", "int", true},
		{"𝓤", "~int", true},
		{"𝓤", "myInt", true},
		{"~int", "int", true},
		{"~int", "~int", true},
		{"~int", "myInt", true},
		{"int", "~int", false},
		{"myInt", "~int", false},
		{"int", "string", false},
		{"~int", "string", false},
		{"int | string", "string", true},
		{"int | string", "~string", false},
		{"~int | string", "int", true},
		{"~int | string", "myInt", true},
		{"~int | string | ∅", "string", true},
		{"~string | ∅ | 𝓤", "myInt", true},
	} {
		xl := maketl(test.xl)
		y := testTerm(test.typ)
		golangt := xl.supersetOf(y)
		if golangt != test.want {
			t.Errorf("(%v).supersetOf(%v) = %v; want %v", test.xl, y, golangt, test.want)
		}
	}
}

func TestTermlistSubsetOf(t *testing.T) {
	for _, test := range []struct {
		xl, yl string
		want   bool
	}{
		{"∅", "∅", true},
		{"∅", "𝓤", true},
		{"𝓤", "∅", false},
		{"𝓤", "𝓤", true},
		{"int", "int | string", true},
		{"~int", "int | string", false},
		{"~int", "myInt | string", false},
		{"myInt", "~int | string", true},
		{"~int", "string | string | int | ~int", true},
		{"myInt", "string | string | ~int", true},
		{"int | string", "string", false},
		{"int | string", "string | int", true},
		{"int | ~string", "string | int", false},
		{"myInt | ~string", "string | int | 𝓤", true},
		{"int | ~string", "string | int | ∅ | string", false},
		{"int | myInt", "string | ~int | ∅ | string", true},
	} {
		xl := maketl(test.xl)
		yl := maketl(test.yl)
		golangt := xl.subsetOf(yl)
		if golangt != test.want {
			t.Errorf("(%v).subsetOf(%v) = %v; want %v", test.xl, test.yl, golangt, test.want)
		}
	}
}
