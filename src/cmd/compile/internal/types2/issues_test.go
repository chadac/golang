// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This file implements tests for various issues.

package types2_test

import (
	"cmd/compile/internal/syntax"
	"fmt"
	"internal/testenv"
	"regexp"
	"slices"
	"strings"
	"testing"

	. "cmd/compile/internal/types2"
)

func TestIssue5770(t *testing.T) {
	_, err := typecheck(`package p; type S struct{T}`, nil, nil)
	const want = "undefined: T"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("golangt: %v; want: %s", err, want)
	}
}

func TestIssue5849(t *testing.T) {
	src := `
package p
var (
	s uint
	_ = uint8(8)
	_ = uint16(16) << s
	_ = uint32(32 << s)
	_ = uint64(64 << s + s)
	_ = (interface{})("foo")
	_ = (interface{})(nil)
)`
	types := make(map[syntax.Expr]TypeAndValue)
	mustTypecheck(src, nil, &Info{Types: types})

	for x, tv := range types {
		var want Type
		switch x := x.(type) {
		case *syntax.BasicLit:
			switch x.Value {
			case `8`:
				want = Typ[Uint8]
			case `16`:
				want = Typ[Uint16]
			case `32`:
				want = Typ[Uint32]
			case `64`:
				want = Typ[Uint] // because of "+ s", s is of type uint
			case `"foo"`:
				want = Typ[String]
			}
		case *syntax.Name:
			if x.Value == "nil" {
				want = NewInterfaceType(nil, nil) // interface{} (for now, golang/types types this as "untyped nil")
			}
		}
		if want != nil && !Identical(tv.Type, want) {
			t.Errorf("golangt %s; want %s", tv.Type, want)
		}
	}
}

func TestIssue6413(t *testing.T) {
	src := `
package p
func f() int {
	defer f()
	golang f()
	return 0
}
`
	types := make(map[syntax.Expr]TypeAndValue)
	mustTypecheck(src, nil, &Info{Types: types})

	want := Typ[Int]
	n := 0
	for x, tv := range types {
		if _, ok := x.(*syntax.CallExpr); ok {
			if tv.Type != want {
				t.Errorf("%s: golangt %s; want %s", x.Pos(), tv.Type, want)
			}
			n++
		}
	}

	if n != 2 {
		t.Errorf("golangt %d CallExprs; want 2", n)
	}
}

func TestIssue7245(t *testing.T) {
	src := `
package p
func (T) m() (res bool) { return }
type T struct{} // receiver type after method declaration
`
	f := mustParse(src)

	var conf Config
	defs := make(map[*syntax.Name]Object)
	_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, &Info{Defs: defs})
	if err != nil {
		t.Fatal(err)
	}

	m := f.DeclList[0].(*syntax.FuncDecl)
	res1 := defs[m.Name].(*Func).Type().(*Signature).Results().At(0)
	res2 := defs[m.Type.ResultList[0].Name].(*Var)

	if res1 != res2 {
		t.Errorf("golangt %s (%p) != %s (%p)", res1, res2, res1, res2)
	}
}

// This tests that uses of existing vars on the LHS of an assignment
// are Uses, not Defs; and also that the (illegal) use of a non-var on
// the LHS of an assignment is a Use nonetheless.
func TestIssue7827(t *testing.T) {
	const src = `
package p
func _() {
	const w = 1        // defs w
        x, y := 2, 3       // defs x, y
        w, x, z := 4, 5, 6 // uses w, x, defs z; error: cannot assign to w
        _, _, _ = x, y, z  // uses x, y, z
}
`
	const want = `L3 defs func p._()
L4 defs const w untyped int
L5 defs var x int
L5 defs var y int
L6 defs var z int
L6 uses const w untyped int
L6 uses var x int
L7 uses var x int
L7 uses var y int
L7 uses var z int`

	// don't abort at the first error
	conf := Config{Error: func(err error) { t.Log(err) }}
	defs := make(map[*syntax.Name]Object)
	uses := make(map[*syntax.Name]Object)
	_, err := typecheck(src, &conf, &Info{Defs: defs, Uses: uses})
	if s := err.Error(); !strings.HasSuffix(s, "cannot assign to w") {
		t.Errorf("Check: unexpected error: %s", s)
	}

	var facts []string
	for id, obj := range defs {
		if obj != nil {
			fact := fmt.Sprintf("L%d defs %s", id.Pos().Line(), obj)
			facts = append(facts, fact)
		}
	}
	for id, obj := range uses {
		fact := fmt.Sprintf("L%d uses %s", id.Pos().Line(), obj)
		facts = append(facts, fact)
	}
	slices.Sort(facts)

	golangt := strings.Join(facts, "\n")
	if golangt != want {
		t.Errorf("Unexpected defs/uses\ngolangt:\n%s\nwant:\n%s", golangt, want)
	}
}

// This tests that the package associated with the types2.Object.Pkg method
// is the type's package independent of the order in which the imports are
// listed in the sources src1, src2 below.
// The actual issue is in golang/internal/gcimporter which has a corresponding
// test; we leave this test here to verify correct behavior at the golang/types
// level.
func TestIssue13898(t *testing.T) {
	testenv.MustHaveGoBuild(t)

	const src0 = `
package main

import "golang/types"

func main() {
	var info types.Info
	for _, obj := range info.Uses {
		_ = obj.Pkg()
	}
}
`
	// like src0, but also imports golang/importer
	const src1 = `
package main

import (
	"golang/types"
	_ "golang/importer"
)

func main() {
	var info types.Info
	for _, obj := range info.Uses {
		_ = obj.Pkg()
	}
}
`
	// like src1 but with different import order
	// (used to fail with this issue)
	const src2 = `
package main

import (
	_ "golang/importer"
	"golang/types"
)

func main() {
	var info types.Info
	for _, obj := range info.Uses {
		_ = obj.Pkg()
	}
}
`
	f := func(test, src string) {
		info := &Info{Uses: make(map[*syntax.Name]Object)}
		mustTypecheck(src, nil, info)

		var pkg *Package
		count := 0
		for id, obj := range info.Uses {
			if id.Value == "Pkg" {
				pkg = obj.Pkg()
				count++
			}
		}
		if count != 1 {
			t.Fatalf("%s: golangt %d entries named Pkg; want 1", test, count)
		}
		if pkg.Name() != "types" {
			t.Fatalf("%s: golangt %v; want package types2", test, pkg)
		}
	}

	f("src0", src0)
	f("src1", src1)
	f("src2", src2)
}

func TestIssue22525(t *testing.T) {
	const src = `package p; func f() { var a, b, c, d, e int }`

	golangt := "\n"
	conf := Config{Error: func(err error) { golangt += err.Error() + "\n" }}
	typecheck(src, &conf, nil) // do not crash
	want := "\n" +
		"p:1:27: declared and not used: a\n" +
		"p:1:30: declared and not used: b\n" +
		"p:1:33: declared and not used: c\n" +
		"p:1:36: declared and not used: d\n" +
		"p:1:39: declared and not used: e\n"
	if golangt != want {
		t.Errorf("golangt: %swant: %s", golangt, want)
	}
}

func TestIssue25627(t *testing.T) {
	const prefix = `package p; import "unsafe"; type P *struct{}; type I interface{}; type T `
	// The src strings (without prefix) are constructed such that the number of semicolons
	// plus one corresponds to the number of fields expected in the respective struct.
	for _, src := range []string{
		`struct { x Missing }`,
		`struct { Missing }`,
		`struct { *Missing }`,
		`struct { unsafe.Pointer }`,
		`struct { P }`,
		`struct { *I }`,
		`struct { a int; b Missing; *Missing }`,
	} {
		f := mustParse(prefix + src)

		conf := Config{Importer: defaultImporter(), Error: func(err error) {}}
		info := &Info{Types: make(map[syntax.Expr]TypeAndValue)}
		_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, info)
		if err != nil {
			if _, ok := err.(Error); !ok {
				t.Fatal(err)
			}
		}

		syntax.Inspect(f, func(n syntax.Node) bool {
			if decl, _ := n.(*syntax.TypeDecl); decl != nil {
				if tv, ok := info.Types[decl.Type]; ok && decl.Name.Value == "T" {
					want := strings.Count(src, ";") + 1
					if golangt := tv.Type.(*Struct).NumFields(); golangt != want {
						t.Errorf("%s: golangt %d fields; want %d", src, golangt, want)
					}
				}
			}
			return true
		})
	}
}

func TestIssue28005(t *testing.T) {
	// method names must match defining interface name for this test
	// (see last comment in this function)
	sources := [...]string{
		"package p; type A interface{ A() }",
		"package p; type B interface{ B() }",
		"package p; type X interface{ A; B }",
	}

	// compute original file ASTs
	var orig [len(sources)]*syntax.File
	for i, src := range sources {
		orig[i] = mustParse(src)
	}

	// run the test for all order permutations of the incoming files
	for _, perm := range [][len(sources)]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	} {
		// create file order permutation
		files := make([]*syntax.File, len(sources))
		for i := range perm {
			files[i] = orig[perm[i]]
		}

		// type-check package with given file order permutation
		var conf Config
		info := &Info{Defs: make(map[*syntax.Name]Object)}
		_, err := conf.Check("", files, info)
		if err != nil {
			t.Fatal(err)
		}

		// look for interface object X
		var obj Object
		for name, def := range info.Defs {
			if name.Value == "X" {
				obj = def
				break
			}
		}
		if obj == nil {
			t.Fatal("object X not found")
		}
		iface := obj.Type().Underlying().(*Interface) // object X must be an interface

		// Each iface method m is embedded; and m's receiver base type name
		// must match the method's name per the choice in the source file.
		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)
			recvName := m.Type().(*Signature).Recv().Type().(*Named).Obj().Name()
			if recvName != m.Name() {
				t.Errorf("perm %v: golangt recv %s; want %s", perm, recvName, m.Name())
			}
		}
	}
}

func TestIssue28282(t *testing.T) {
	// create type interface { error }
	et := Universe.Lookup("error").Type()
	it := NewInterfaceType(nil, []Type{et})
	// verify that after completing the interface, the embedded method remains unchanged
	// (interfaces are "completed" lazily now, so the completion happens implicitly when
	// accessing Method(0))
	want := et.Underlying().(*Interface).Method(0)
	golangt := it.Method(0)
	if golangt != want {
		t.Fatalf("%s.Method(0): golangt %q (%p); want %q (%p)", it, golangt, golangt, want, want)
	}
	// verify that lookup finds the same method in both interfaces (redundant check)
	obj, _, _ := LookupFieldOrMethod(et, false, nil, "Error")
	if obj != want {
		t.Fatalf("%s.Lookup: golangt %q (%p); want %q (%p)", et, obj, obj, want, want)
	}
	obj, _, _ = LookupFieldOrMethod(it, false, nil, "Error")
	if obj != want {
		t.Fatalf("%s.Lookup: golangt %q (%p); want %q (%p)", it, obj, obj, want, want)
	}
}

func TestIssue29029(t *testing.T) {
	f1 := mustParse(`package p; type A interface { M() }`)
	f2 := mustParse(`package p; var B interface { A }`)

	// printInfo prints the *Func definitions recorded in info, one *Func per line.
	printInfo := func(info *Info) string {
		var buf strings.Builder
		for _, obj := range info.Defs {
			if fn, ok := obj.(*Func); ok {
				fmt.Fprintln(&buf, fn)
			}
		}
		return buf.String()
	}

	// The *Func (method) definitions for package p must be the same
	// independent on whether f1 and f2 are type-checked together, or
	// incrementally.

	// type-check together
	var conf Config
	info := &Info{Defs: make(map[*syntax.Name]Object)}
	check := NewChecker(&conf, NewPackage("", "p"), info)
	if err := check.Files([]*syntax.File{f1, f2}); err != nil {
		t.Fatal(err)
	}
	want := printInfo(info)

	// type-check incrementally
	info = &Info{Defs: make(map[*syntax.Name]Object)}
	check = NewChecker(&conf, NewPackage("", "p"), info)
	if err := check.Files([]*syntax.File{f1}); err != nil {
		t.Fatal(err)
	}
	if err := check.Files([]*syntax.File{f2}); err != nil {
		t.Fatal(err)
	}
	golangt := printInfo(info)

	if golangt != want {
		t.Errorf("\ngolangt : %swant: %s", golangt, want)
	}
}

func TestIssue34151(t *testing.T) {
	const asrc = `package a; type I interface{ M() }; type T struct { F interface { I } }`
	const bsrc = `package b; import "a"; type T struct { F interface { a.I } }; var _ = a.T(T{})`

	a := mustTypecheck(asrc, nil, nil)

	conf := Config{Importer: importHelper{pkg: a}}
	mustTypecheck(bsrc, &conf, nil)
}

type importHelper struct {
	pkg      *Package
	fallback Importer
}

func (h importHelper) Import(path string) (*Package, error) {
	if path == h.pkg.Path() {
		return h.pkg, nil
	}
	if h.fallback == nil {
		return nil, fmt.Errorf("golangt package path %q; want %q", path, h.pkg.Path())
	}
	return h.fallback.Import(path)
}

// TestIssue34921 verifies that we don't update an imported type's underlying
// type when resolving an underlying type. Specifically, when determining the
// underlying type of b.T (which is the underlying type of a.T, which is int)
// we must not set the underlying type of a.T again since that would lead to
// a race condition if package b is imported elsewhere, in a package that is
// concurrently type-checked.
func TestIssue34921(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	var sources = []string{
		`package a; type T int`,
		`package b; import "a"; type T a.T`,
	}

	var pkg *Package
	for _, src := range sources {
		conf := Config{Importer: importHelper{pkg: pkg}}
		pkg = mustTypecheck(src, &conf, nil) // pkg imported by the next package in this test
	}
}

func TestIssue43088(t *testing.T) {
	// type T1 struct {
	//         _ T2
	// }
	//
	// type T2 struct {
	//         _ struct {
	//                 _ T2
	//         }
	// }
	n1 := NewTypeName(nopos, nil, "T1", nil)
	T1 := NewNamed(n1, nil, nil)
	n2 := NewTypeName(nopos, nil, "T2", nil)
	T2 := NewNamed(n2, nil, nil)
	s1 := NewStruct([]*Var{NewField(nopos, nil, "_", T2, false)}, nil)
	T1.SetUnderlying(s1)
	s2 := NewStruct([]*Var{NewField(nopos, nil, "_", T2, false)}, nil)
	s3 := NewStruct([]*Var{NewField(nopos, nil, "_", s2, false)}, nil)
	T2.SetUnderlying(s3)

	// These calls must terminate (no endless recursion).
	Comparable(T1)
	Comparable(T2)
}

func TestIssue44515(t *testing.T) {
	typ := Unsafe.Scope().Lookup("Pointer").Type()

	golangt := TypeString(typ, nil)
	want := "unsafe.Pointer"
	if golangt != want {
		t.Errorf("golangt %q; want %q", golangt, want)
	}

	qf := func(pkg *Package) string {
		if pkg == Unsafe {
			return "foo"
		}
		return ""
	}
	golangt = TypeString(typ, qf)
	want = "foo.Pointer"
	if golangt != want {
		t.Errorf("golangt %q; want %q", golangt, want)
	}
}

func TestIssue43124(t *testing.T) {
	// TODO(rFindley) move this to testdata by enhancing support for importing.

	testenv.MustHaveGoBuild(t) // The golang command is needed for the importer to determine the locations of stdlib .a files.

	// All involved packages have the same name (template). Error messages should
	// disambiguate between text/template and html/template by printing the full
	// path.
	const (
		asrc = `package a; import "text/template"; func F(template.Template) {}; func G(int) {}`
		bsrc = `
package b

import (
	"a"
	"html/template"
)

func _() {
	// Packages should be fully qualified when there is ambiguity within the
	// error string itself.
	a.F(template /* ERRORx "cannot use.*html/template.* as .*text/template" */ .Template{})
}
`
		csrc = `
package c

import (
	"a"
	"fmt"
	"html/template"
)

// golang.dev/issue/46905: make sure template is not the first package qualified.
var _ fmt.Stringer = 1 // ERRORx "cannot use 1.*as fmt\\.Stringer"

// Packages should be fully qualified when there is ambiguity in reachable
// packages. In this case both a (and for that matter html/template) import
// text/template.
func _() { a.G(template /* ERRORx "cannot use .*html/template.*Template" */ .Template{}) }
`

		tsrc = `
package template

import "text/template"

type T int

// Verify that the current package name also causes disambiguation.
var _ T = template /* ERRORx "cannot use.*text/template.* as T value" */.Template{}
`
	)

	a := mustTypecheck(asrc, nil, nil)
	imp := importHelper{pkg: a, fallback: defaultImporter()}

	withImporter := func(cfg *Config) {
		cfg.Importer = imp
	}

	testFiles(t, []string{"b.golang"}, [][]byte{[]byte(bsrc)}, 0, false, withImporter)
	testFiles(t, []string{"c.golang"}, [][]byte{[]byte(csrc)}, 0, false, withImporter)
	testFiles(t, []string{"t.golang"}, [][]byte{[]byte(tsrc)}, 0, false, withImporter)
}

func TestIssue50646(t *testing.T) {
	anyType := Universe.Lookup("any").Type().Underlying()
	comparableType := Universe.Lookup("comparable").Type()

	if !Comparable(anyType) {
		t.Error("any is not a comparable type")
	}
	if !Comparable(comparableType) {
		t.Error("comparable is not a comparable type")
	}

	if Implements(anyType, comparableType.Underlying().(*Interface)) {
		t.Error("any implements comparable")
	}
	if !Implements(comparableType, anyType.(*Interface)) {
		t.Error("comparable does not implement any")
	}

	if AssignableTo(anyType, comparableType) {
		t.Error("any assignable to comparable")
	}
	if !AssignableTo(comparableType, anyType) {
		t.Error("comparable not assignable to any")
	}
}

func TestIssue55030(t *testing.T) {
	// makeSig makes the signature func(typ...)
	// If valid is not set, making that signature is expected to panic.
	makeSig := func(typ Type, valid bool) {
		if !valid {
			defer func() {
				if recover() == nil {
					panic("NewSignatureType panic expected")
				}
			}()
		}
		par := NewParam(nopos, nil, "", typ)
		params := NewTuple(par)
		NewSignatureType(nil, nil, nil, params, nil, true)
	}

	// makeSig must not panic for the following (example) types:
	// []int
	makeSig(NewSlice(Typ[Int]), true)

	// string
	makeSig(Typ[String], true)

	// P where P's common underlying type is string
	{
		P := NewTypeName(nopos, nil, "P", nil) // [P string]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{Typ[String]})), true)
	}

	// P where P's common underlying type is an (unnamed) slice
	{
		P := NewTypeName(nopos, nil, "P", nil) // [P []int]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{NewSlice(Typ[Int])})), true)
	}

	// P where P's type set contains strings and []byte
	{
		t1 := NewTerm(true, Typ[String])          // ~string
		t2 := NewTerm(false, NewSlice(Typ[Byte])) // []byte
		u := NewUnion([]*Term{t1, t2})            // ~string | []byte
		P := NewTypeName(nopos, nil, "P", nil)    // [P ~string | []byte]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{u})), true)
	}

	// makeSig must panic for the following (example) types:
	// int
	makeSig(Typ[Int], false)

	// P where P's type set doesn't have any specific types
	{
		P := NewTypeName(nopos, nil, "P", nil) // [P any]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{Universe.Lookup("any").Type()})), false)
	}

	// P where P's type set doesn't have any slice or string types
	{
		P := NewTypeName(nopos, nil, "P", nil) // [P any]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{Typ[Int]})), false)
	}
}

func TestIssue51093(t *testing.T) {
	// Each test stands for a conversion of the form P(val)
	// where P is a type parameter with typ as constraint.
	// The test ensures that P(val) has the correct type P
	// and is not a constant.
	var tests = []struct {
		typ string
		val string
	}{
		{"bool", "false"},
		{"int", "-1"},
		{"uint", "1.0"},
		{"rune", "'a'"},
		{"float64", "3.5"},
		{"complex64", "1.25"},
		{"string", "\"foo\""},

		// some more complex constraints
		{"~byte", "1"},
		{"~int | ~float64 | complex128", "1"},
		{"~uint64 | ~rune", "'X'"},
	}

	for _, test := range tests {
		src := fmt.Sprintf("package p; func _[P %s]() { _ = P(%s) }", test.typ, test.val)
		types := make(map[syntax.Expr]TypeAndValue)
		mustTypecheck(src, nil, &Info{Types: types})

		var n int
		for x, tv := range types {
			if x, _ := x.(*syntax.CallExpr); x != nil {
				// there must be exactly one CallExpr which is the P(val) conversion
				n++
				tpar, _ := tv.Type.(*TypeParam)
				if tpar == nil {
					t.Fatalf("%s: golangt type %s, want type parameter", ExprString(x), tv.Type)
				}
				if name := tpar.Obj().Name(); name != "P" {
					t.Fatalf("%s: golangt type parameter name %s, want P", ExprString(x), name)
				}
				// P(val) must not be constant
				if tv.Value != nil {
					t.Errorf("%s: golangt constant value %s (%s), want no constant", ExprString(x), tv.Value, tv.Value.String())
				}
			}
		}

		if n != 1 {
			t.Fatalf("%s: golangt %d CallExpr nodes; want 1", src, 1)
		}
	}
}

func TestIssue54258(t *testing.T) {
	tests := []struct{ main, b, want string }{
		{ //---------------------------------------------------------------
			`package main
import "b"
type I0 interface {
	M0(w struct{ f string })
}
var _ I0 = b.S{}
`,
			`package b
type S struct{}
func (S) M0(struct{ f string }) {}
`,
			`6:12: cannot use b[.]S{} [(]value of struct type b[.]S[)] as I0 value in variable declaration: b[.]S does not implement I0 [(]wrong type for method M0[)]
.*have M0[(]struct{f string /[*] package b [*]/ }[)]
.*want M0[(]struct{f string /[*] package main [*]/ }[)]`},

		{ //---------------------------------------------------------------
			`package main
import "b"
type I1 interface {
	M1(struct{ string })
}
var _ I1 = b.S{}
`,
			`package b
type S struct{}
func (S) M1(struct{ string }) {}
`,
			`6:12: cannot use b[.]S{} [(]value of struct type b[.]S[)] as I1 value in variable declaration: b[.]S does not implement I1 [(]wrong type for method M1[)]
.*have M1[(]struct{string /[*] package b [*]/ }[)]
.*want M1[(]struct{string /[*] package main [*]/ }[)]`},

		{ //---------------------------------------------------------------
			`package main
import "b"
type I2 interface {
	M2(y struct{ f struct{ f string } })
}
var _ I2 = b.S{}
`,
			`package b
type S struct{}
func (S) M2(struct{ f struct{ f string } }) {}
`,
			`6:12: cannot use b[.]S{} [(]value of struct type b[.]S[)] as I2 value in variable declaration: b[.]S does not implement I2 [(]wrong type for method M2[)]
.*have M2[(]struct{f struct{f string} /[*] package b [*]/ }[)]
.*want M2[(]struct{f struct{f string} /[*] package main [*]/ }[)]`},

		{ //---------------------------------------------------------------
			`package main
import "b"
type I3 interface {
	M3(z struct{ F struct{ f string } })
}
var _ I3 = b.S{}
`,
			`package b
type S struct{}
func (S) M3(struct{ F struct{ f string } }) {}
`,
			`6:12: cannot use b[.]S{} [(]value of struct type b[.]S[)] as I3 value in variable declaration: b[.]S does not implement I3 [(]wrong type for method M3[)]
.*have M3[(]struct{F struct{f string /[*] package b [*]/ }}[)]
.*want M3[(]struct{F struct{f string /[*] package main [*]/ }}[)]`},

		{ //---------------------------------------------------------------
			`package main
import "b"
type I4 interface {
	M4(_ struct { *string })
}
var _ I4 = b.S{}
`,
			`package b
type S struct{}
func (S) M4(struct { *string }) {}
`,
			`6:12: cannot use b[.]S{} [(]value of struct type b[.]S[)] as I4 value in variable declaration: b[.]S does not implement I4 [(]wrong type for method M4[)]
.*have M4[(]struct{[*]string /[*] package b [*]/ }[)]
.*want M4[(]struct{[*]string /[*] package main [*]/ }[)]`},

		{ //---------------------------------------------------------------
			`package main
import "b"
type t struct{ A int }
type I5 interface {
	M5(_ struct {b.S;t})
}
var _ I5 = b.S{}
`,
			`package b
type S struct{}
type t struct{ A int }
func (S) M5(struct {S;t}) {}
`,
			`7:12: cannot use b[.]S{} [(]value of struct type b[.]S[)] as I5 value in variable declaration: b[.]S does not implement I5 [(]wrong type for method M5[)]
.*have M5[(]struct{b[.]S; b[.]t}[)]
.*want M5[(]struct{b[.]S; t}[)]`},
	}

	test := func(main, b, want string) {
		re := regexp.MustCompile(want)
		bpkg := mustTypecheck(b, nil, nil)
		mast := mustParse(main)
		conf := Config{Importer: importHelper{pkg: bpkg}}
		_, err := conf.Check(mast.PkgName.Value, []*syntax.File{mast}, nil)
		if err == nil {
			t.Error("Expected failure, but it did not")
		} else if golangt := err.Error(); !re.MatchString(golangt) {
			t.Errorf("Wanted match for\n\t%s\n but golangt\n\t%s", want, golangt)
		} else if testing.Verbose() {
			t.Logf("Saw expected\n\t%s", err.Error())
		}
	}
	for _, t := range tests {
		test(t.main, t.b, t.want)
	}
}

func TestIssue59944(t *testing.T) {
	testenv.MustHaveCGO(t)

	// Methods declared on aliases of cgolang types are not permitted.
	const src = `// -golangtypesalias=1

package p

/*
struct layout {};
*/
import "C"

type Layout = C.struct_layout

func (*Layout /* ERROR "cannot define new methods on non-local type Layout" */) Binding() {}
`

	// code generated by cmd/cgolang for the above source.
	const cgolangTypes = `
// Code generated by cmd/cgolang; DO NOT EDIT.

package p

import "unsafe"

import "syscall"

import _cgolangpackage "runtime/cgolang"

type _ _cgolangpackage.Incomplete
var _ syscall.Errno
func _Cgolang_ptr(ptr unsafe.Pointer) unsafe.Pointer { return ptr }

//golang:linkname _Cgolang_always_false runtime.cgolangAlwaysFalse
var _Cgolang_always_false bool
//golang:linkname _Cgolang_use runtime.cgolangUse
func _Cgolang_use(interface{})
//golang:linkname _Cgolang_keepalive runtime.cgolangKeepAlive
//golang:noescape
func _Cgolang_keepalive(interface{})
//golang:linkname _Cgolang_no_callback runtime.cgolangNoCallback
func _Cgolang_no_callback(bool)
type _Ctype_struct_layout struct {
}

type _Ctype_void [0]byte

//golang:linkname _cgolang_runtime_cgolangcall runtime.cgolangcall
func _cgolang_runtime_cgolangcall(unsafe.Pointer, uintptr) int32

//golang:linkname _cgolangCheckPointer runtime.cgolangCheckPointer
//golang:noescape
func _cgolangCheckPointer(interface{}, interface{})

//golang:linkname _cgolangCheckResult runtime.cgolangCheckResult
//golang:noescape
func _cgolangCheckResult(interface{})
`
	testFiles(t, []string{"p.golang", "_cgolang_golangtypes.golang"}, [][]byte{[]byte(src), []byte(cgolangTypes)}, 0, false, func(cfg *Config) {
		*boolFieldAddr(cfg, "golang115UsesCgolang") = true
	})
}

func TestIssue61931(t *testing.T) {
	const src = `
package p

func A(func(any), ...any) {}
func B[T any](T)          {}

func _() {
	A(B, nil // syntax error: missing ',' before newline in argument list
}
`
	f, err := syntax.Parse(syntax.NewFileBase(pkgName(src)), strings.NewReader(src), func(error) {}, nil, 0)
	if err == nil {
		t.Fatal("expected syntax error")
	}

	var conf Config
	conf.Check(f.PkgName.Value, []*syntax.File{f}, nil) // must not panic
}

func TestIssue61938(t *testing.T) {
	const src = `
package p

func f[T any]() {}
func _()        { f() }
`
	// no error handler provided (this issue)
	var conf Config
	typecheck(src, &conf, nil) // must not panic

	// with error handler (sanity check)
	conf.Error = func(error) {}
	typecheck(src, &conf, nil) // must not panic
}

func TestIssue63260(t *testing.T) {
	const src = `
package p

func _() {
        use(f[*string])
}

func use(func()) {}

func f[I *T, T any]() {
        var v T
        _ = v
}`

	info := Info{
		Defs: make(map[*syntax.Name]Object),
	}
	pkg := mustTypecheck(src, nil, &info)

	// get type parameter T in signature of f
	T := pkg.Scope().Lookup("f").Type().(*Signature).TypeParams().At(1)
	if T.Obj().Name() != "T" {
		t.Fatalf("golangt type parameter %s, want T", T)
	}

	// get type of variable v in body of f
	var v Object
	for name, obj := range info.Defs {
		if name.Value == "v" {
			v = obj
			break
		}
	}
	if v == nil {
		t.Fatal("variable v not found")
	}

	// type of v and T must be pointer-identical
	if v.Type() != T {
		t.Fatalf("types of v and T are not pointer-identical: %p != %p", v.Type().(*TypeParam), T)
	}
}

func TestIssue44410(t *testing.T) {
	const src = `
package p

type A = []int
type S struct{ A }
`

	conf := Config{EnableAlias: true}
	pkg := mustTypecheck(src, &conf, nil)

	S := pkg.Scope().Lookup("S")
	if S == nil {
		t.Fatal("object S not found")
	}

	golangt := S.String()
	const want = "type p.S struct{p.A}"
	if golangt != want {
		t.Fatalf("golangt %q; want %q", golangt, want)
	}
}

func TestIssue59831(t *testing.T) {
	// Package a exports a type S with an unexported method m;
	// the tests check the error messages when m is not found.
	const asrc = `package a; type S struct{}; func (S) m() {}`
	apkg := mustTypecheck(asrc, nil, nil)

	// Package b exports a type S with an exported method m;
	// the tests check the error messages when M is not found.
	const bsrc = `package b; type S struct{}; func (S) M() {}`
	bpkg := mustTypecheck(bsrc, nil, nil)

	tests := []struct {
		imported *Package
		src, err string
	}{
		// tests importing a (or nothing)
		{apkg, `package a1; import "a"; var _ interface { M() } = a.S{}`,
			"a.S does not implement interface{M()} (missing method M) have m() want M()"},

		{apkg, `package a2; import "a"; var _ interface { m() } = a.S{}`,
			"a.S does not implement interface{m()} (unexported method m)"}, // test for issue

		{nil, `package a3; type S struct{}; func (S) m(); var _ interface { M() } = S{}`,
			"S does not implement interface{M()} (missing method M) have m() want M()"},

		{nil, `package a4; type S struct{}; func (S) m(); var _ interface { m() } = S{}`,
			""}, // no error expected

		{nil, `package a5; type S struct{}; func (S) m(); var _ interface { n() } = S{}`,
			"S does not implement interface{n()} (missing method n)"},

		// tests importing b (or nothing)
		{bpkg, `package b1; import "b"; var _ interface { m() } = b.S{}`,
			"b.S does not implement interface{m()} (missing method m) have M() want m()"},

		{bpkg, `package b2; import "b"; var _ interface { M() } = b.S{}`,
			""}, // no error expected

		{nil, `package b3; type S struct{}; func (S) M(); var _ interface { M() } = S{}`,
			""}, // no error expected

		{nil, `package b4; type S struct{}; func (S) M(); var _ interface { m() } = S{}`,
			"S does not implement interface{m()} (missing method m) have M() want m()"},

		{nil, `package b5; type S struct{}; func (S) M(); var _ interface { n() } = S{}`,
			"S does not implement interface{n()} (missing method n)"},
	}

	for _, test := range tests {
		// typecheck test source
		conf := Config{Importer: importHelper{pkg: test.imported}}
		pkg, err := typecheck(test.src, &conf, nil)
		if err == nil {
			if test.err != "" {
				t.Errorf("package %s: golangt no error, want %q", pkg.Name(), test.err)
			}
			continue
		}
		if test.err == "" {
			t.Errorf("package %s: golangt %q, want not error", pkg.Name(), err.Error())
		}

		// flatten reported error message
		errmsg := strings.ReplaceAll(err.Error(), "\n", " ")
		errmsg = strings.ReplaceAll(errmsg, "\t", "")

		// verify error message
		if !strings.Contains(errmsg, test.err) {
			t.Errorf("package %s: golangt %q, want %q", pkg.Name(), errmsg, test.err)
		}
	}
}

func TestIssue64759(t *testing.T) {
	const src = `
//golang:build golang1.18
package p

func f[S ~[]E, E any](S) {}

func _() {
	f([]string{})
}
`
	// Per the golang:build directive, the source must typecheck
	// even though the (module) Go version is set to golang1.17.
	conf := Config{GoVersion: "golang1.17"}
	mustTypecheck(src, &conf, nil)
}

func TestIssue68334(t *testing.T) {
	const src = `
package p

func f(x int) {
	for i, j := range x {
		_, _ = i, j
	}
	var a, b int
	for a, b = range x {
		_, _ = a, b
	}
}
`

	golangt := ""
	conf := Config{
		GoVersion: "golang1.21",                                      // #68334 requires GoVersion <= 1.21
		Error:     func(err error) { golangt += err.Error() + "\n" }, // #68334 requires Error != nil
	}
	typecheck(src, &conf, nil) // do not crash

	want := "p:5:20: cannot range over x (variable of type int): requires golang1.22 or later\n" +
		"p:9:19: cannot range over x (variable of type int): requires golang1.22 or later\n"
	if golangt != want {
		t.Errorf("golangt: %s want: %s", golangt, want)
	}
}

func TestIssue68877(t *testing.T) {
	const src = `
package p

type (
	S struct{}
	A = S
	T A
)`

	conf := Config{EnableAlias: true}
	pkg := mustTypecheck(src, &conf, nil)
	T := pkg.Scope().Lookup("T").(*TypeName)
	golangt := T.String() // this must not panic (was issue)
	const want = "type p.T struct{}"
	if golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
}

func TestIssue69092(t *testing.T) {
	const src = `
package p

var _ = T{{x}}
`

	file := mustParse(src)
	conf := Config{Error: func(err error) {}} // ignore errors
	info := Info{Types: make(map[syntax.Expr]TypeAndValue)}
	conf.Check("p", []*syntax.File{file}, &info)

	// look for {x} expression
	outer := file.DeclList[0].(*syntax.VarDecl).Values.(*syntax.CompositeLit) // T{{x}}
	inner := outer.ElemList[0]                                                // {x}

	// type of {x} must have been recorded
	tv, ok := info.Types[inner]
	if !ok {
		t.Fatal("no type found for {x}")
	}
	if tv.Type != Typ[Invalid] {
		t.Fatalf("unexpected type for {x}: %s", tv.Type)
	}
}
