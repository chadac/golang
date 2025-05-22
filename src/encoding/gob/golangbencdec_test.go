// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This file contains tests of the GolangbEncoder/GolangbDecoder support.

package golangb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Types that implement the GolangbEncoder/Decoder interfaces.

type ByteStruct struct {
	a byte // not an exported field
}

type StringStruct struct {
	s string // not an exported field
}

type ArrayStruct struct {
	a [8192]byte // not an exported field
}

type Golangbber int

type ValueGolangbber string // encodes with a value, decodes with a pointer.

type BinaryGolangbber int

type BinaryValueGolangbber string

type TextGolangbber int

type TextValueGolangbber string

// The relevant methods

func (g *ByteStruct) GolangbEncode() ([]byte, error) {
	b := make([]byte, 3)
	b[0] = g.a
	b[1] = g.a + 1
	b[2] = g.a + 2
	return b, nil
}

func (g *ByteStruct) GolangbDecode(data []byte) error {
	if g == nil {
		return errors.New("NIL RECEIVER")
	}
	// Expect N sequential-valued bytes.
	if len(data) == 0 {
		return io.EOF
	}
	g.a = data[0]
	for i, c := range data {
		if c != g.a+byte(i) {
			return errors.New("invalid data sequence")
		}
	}
	return nil
}

func (g *StringStruct) GolangbEncode() ([]byte, error) {
	return []byte(g.s), nil
}

func (g *StringStruct) GolangbDecode(data []byte) error {
	// Expect N sequential-valued bytes.
	if len(data) == 0 {
		return io.EOF
	}
	a := data[0]
	for i, c := range data {
		if c != a+byte(i) {
			return errors.New("invalid data sequence")
		}
	}
	g.s = string(data)
	return nil
}

func (a *ArrayStruct) GolangbEncode() ([]byte, error) {
	return a.a[:], nil
}

func (a *ArrayStruct) GolangbDecode(data []byte) error {
	if len(data) != len(a.a) {
		return errors.New("wrong length in array decode")
	}
	copy(a.a[:], data)
	return nil
}

func (g *Golangbber) GolangbEncode() ([]byte, error) {
	return []byte(fmt.Sprintf("VALUE=%d", *g)), nil
}

func (g *Golangbber) GolangbDecode(data []byte) error {
	_, err := fmt.Sscanf(string(data), "VALUE=%d", (*int)(g))
	return err
}

func (g *BinaryGolangbber) MarshalBinary() ([]byte, error) {
	return []byte(fmt.Sprintf("VALUE=%d", *g)), nil
}

func (g *BinaryGolangbber) UnmarshalBinary(data []byte) error {
	_, err := fmt.Sscanf(string(data), "VALUE=%d", (*int)(g))
	return err
}

func (g *TextGolangbber) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("VALUE=%d", *g)), nil
}

func (g *TextGolangbber) UnmarshalText(data []byte) error {
	_, err := fmt.Sscanf(string(data), "VALUE=%d", (*int)(g))
	return err
}

func (v ValueGolangbber) GolangbEncode() ([]byte, error) {
	return []byte(fmt.Sprintf("VALUE=%s", v)), nil
}

func (v *ValueGolangbber) GolangbDecode(data []byte) error {
	_, err := fmt.Sscanf(string(data), "VALUE=%s", (*string)(v))
	return err
}

func (v BinaryValueGolangbber) MarshalBinary() ([]byte, error) {
	return []byte(fmt.Sprintf("VALUE=%s", v)), nil
}

func (v *BinaryValueGolangbber) UnmarshalBinary(data []byte) error {
	_, err := fmt.Sscanf(string(data), "VALUE=%s", (*string)(v))
	return err
}

func (v TextValueGolangbber) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("VALUE=%s", v)), nil
}

func (v *TextValueGolangbber) UnmarshalText(data []byte) error {
	_, err := fmt.Sscanf(string(data), "VALUE=%s", (*string)(v))
	return err
}

// Structs that include GolangbEncodable fields.

type GolangbTest0 struct {
	X int // guarantee we have  something in common with GolangbTest*
	G *ByteStruct
}

type GolangbTest1 struct {
	X int // guarantee we have  something in common with GolangbTest*
	G *StringStruct
}

type GolangbTest2 struct {
	X int    // guarantee we have  something in common with GolangbTest*
	G string // not a GolangbEncoder - should give us errors
}

type GolangbTest3 struct {
	X int // guarantee we have  something in common with GolangbTest*
	G *Golangbber
	B *BinaryGolangbber
	T *TextGolangbber
}

type GolangbTest4 struct {
	X  int // guarantee we have  something in common with GolangbTest*
	V  ValueGolangbber
	BV BinaryValueGolangbber
	TV TextValueGolangbber
}

type GolangbTest5 struct {
	X  int // guarantee we have  something in common with GolangbTest*
	V  *ValueGolangbber
	BV *BinaryValueGolangbber
	TV *TextValueGolangbber
}

type GolangbTest6 struct {
	X  int // guarantee we have  something in common with GolangbTest*
	V  ValueGolangbber
	W  *ValueGolangbber
	BV BinaryValueGolangbber
	BW *BinaryValueGolangbber
	TV TextValueGolangbber
	TW *TextValueGolangbber
}

type GolangbTest7 struct {
	X  int // guarantee we have  something in common with GolangbTest*
	V  *ValueGolangbber
	W  ValueGolangbber
	BV *BinaryValueGolangbber
	BW BinaryValueGolangbber
	TV *TextValueGolangbber
	TW TextValueGolangbber
}

type GolangbTestIgnoreEncoder struct {
	X int // guarantee we have  something in common with GolangbTest*
}

type GolangbTestValueEncDec struct {
	X int          // guarantee we have  something in common with GolangbTest*
	G StringStruct // not a pointer.
}

type GolangbTestIndirectEncDec struct {
	X int             // guarantee we have  something in common with GolangbTest*
	G ***StringStruct // indirections to the receiver.
}

type GolangbTestArrayEncDec struct {
	X int         // guarantee we have  something in common with GolangbTest*
	A ArrayStruct // not a pointer.
}

type GolangbTestIndirectArrayEncDec struct {
	X int            // guarantee we have  something in common with GolangbTest*
	A ***ArrayStruct // indirections to a large receiver.
}

func TestGolangbEncoderField(t *testing.T) {
	b := new(bytes.Buffer)
	// First a field that's a structure.
	enc := NewEncoder(b)
	err := enc.Encode(GolangbTest0{17, &ByteStruct{'A'}})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTest0)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x.G.a != 'A' {
		t.Errorf("expected 'A' golangt %c", x.G.a)
	}
	// Now a field that's not a structure.
	b.Reset()
	golangbber := Golangbber(23)
	bgolangbber := BinaryGolangbber(24)
	tgolangbber := TextGolangbber(25)
	err = enc.Encode(GolangbTest3{17, &golangbber, &bgolangbber, &tgolangbber})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	y := new(GolangbTest3)
	err = dec.Decode(y)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if *y.G != 23 || *y.B != 24 || *y.T != 25 {
		t.Errorf("expected '23 golangt %d", *y.G)
	}
}

// Even though the field is a value, we can still take its address
// and should be able to call the methods.
func TestGolangbEncoderValueField(t *testing.T) {
	b := new(bytes.Buffer)
	// First a field that's a structure.
	enc := NewEncoder(b)
	err := enc.Encode(&GolangbTestValueEncDec{17, StringStruct{"HIJKL"}})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTestValueEncDec)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x.G.s != "HIJKL" {
		t.Errorf("expected `HIJKL` golangt %s", x.G.s)
	}
}

// GolangbEncode/Decode should work even if the value is
// more indirect than the receiver.
func TestGolangbEncoderIndirectField(t *testing.T) {
	b := new(bytes.Buffer)
	// First a field that's a structure.
	enc := NewEncoder(b)
	s := &StringStruct{"HIJKL"}
	sp := &s
	err := enc.Encode(GolangbTestIndirectEncDec{17, &sp})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTestIndirectEncDec)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if (***x.G).s != "HIJKL" {
		t.Errorf("expected `HIJKL` golangt %s", (***x.G).s)
	}
}

// Test with a large field with methods.
func TestGolangbEncoderArrayField(t *testing.T) {
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	var a GolangbTestArrayEncDec
	a.X = 17
	for i := range a.A.a {
		a.A.a[i] = byte(i)
	}
	err := enc.Encode(&a)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTestArrayEncDec)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	for i, v := range x.A.a {
		if v != byte(i) {
			t.Errorf("expected %x golangt %x", byte(i), v)
			break
		}
	}
}

// Test an indirection to a large field with methods.
func TestGolangbEncoderIndirectArrayField(t *testing.T) {
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	var a GolangbTestIndirectArrayEncDec
	a.X = 17
	var array ArrayStruct
	ap := &array
	app := &ap
	a.A = &app
	for i := range array.a {
		array.a[i] = byte(i)
	}
	err := enc.Encode(a)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTestIndirectArrayEncDec)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	for i, v := range (***x.A).a {
		if v != byte(i) {
			t.Errorf("expected %x golangt %x", byte(i), v)
			break
		}
	}
}

// As long as the fields have the same name and implement the
// interface, we can cross-connect them. Not sure it's useful
// and may even be bad but it works and it's hard to prevent
// without exposing the contents of the object, which would
// defeat the purpose.
func TestGolangbEncoderFieldsOfDifferentType(t *testing.T) {
	// first, string in field to byte in field
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	err := enc.Encode(GolangbTest1{17, &StringStruct{"ABC"}})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTest0)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x.G.a != 'A' {
		t.Errorf("expected 'A' golangt %c", x.G.a)
	}
	// now the other direction, byte in field to string in field
	b.Reset()
	err = enc.Encode(GolangbTest0{17, &ByteStruct{'X'}})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	y := new(GolangbTest1)
	err = dec.Decode(y)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if y.G.s != "XYZ" {
		t.Fatalf("expected `XYZ` golangt %q", y.G.s)
	}
}

// Test that we can encode a value and decode into a pointer.
func TestGolangbEncoderValueEncoder(t *testing.T) {
	// first, string in field to byte in field
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	err := enc.Encode(GolangbTest4{17, ValueGolangbber("hello"), BinaryValueGolangbber("Καλημέρα"), TextValueGolangbber("こんにちは")})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTest5)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if *x.V != "hello" || *x.BV != "Καλημέρα" || *x.TV != "こんにちは" {
		t.Errorf("expected `hello` golangt %s", *x.V)
	}
}

// Test that we can use a value then a pointer type of a GolangbEncoder
// in the same encoded value. Bug 4647.
func TestGolangbEncoderValueThenPointer(t *testing.T) {
	v := ValueGolangbber("forty-two")
	w := ValueGolangbber("six-by-nine")
	bv := BinaryValueGolangbber("1nanocentury")
	bw := BinaryValueGolangbber("πseconds")
	tv := TextValueGolangbber("gravitationalacceleration")
	tw := TextValueGolangbber("π²ft/s²")

	// this was a bug: encoding a GolangbEncoder by value before a GolangbEncoder
	// pointer would cause duplicate type definitions to be sent.

	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	if err := enc.Encode(GolangbTest6{42, v, &w, bv, &bw, tv, &tw}); err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTest6)
	if err := dec.Decode(x); err != nil {
		t.Fatal("decode error:", err)
	}

	if golangt, want := x.V, v; golangt != want {
		t.Errorf("v = %q, want %q", golangt, want)
	}
	if golangt, want := x.W, w; golangt == nil {
		t.Errorf("w = nil, want %q", want)
	} else if *golangt != want {
		t.Errorf("w = %q, want %q", *golangt, want)
	}

	if golangt, want := x.BV, bv; golangt != want {
		t.Errorf("bv = %q, want %q", golangt, want)
	}
	if golangt, want := x.BW, bw; golangt == nil {
		t.Errorf("bw = nil, want %q", want)
	} else if *golangt != want {
		t.Errorf("bw = %q, want %q", *golangt, want)
	}

	if golangt, want := x.TV, tv; golangt != want {
		t.Errorf("tv = %q, want %q", golangt, want)
	}
	if golangt, want := x.TW, tw; golangt == nil {
		t.Errorf("tw = nil, want %q", want)
	} else if *golangt != want {
		t.Errorf("tw = %q, want %q", *golangt, want)
	}
}

// Test that we can use a pointer then a value type of a GolangbEncoder
// in the same encoded value.
func TestGolangbEncoderPointerThenValue(t *testing.T) {
	v := ValueGolangbber("forty-two")
	w := ValueGolangbber("six-by-nine")
	bv := BinaryValueGolangbber("1nanocentury")
	bw := BinaryValueGolangbber("πseconds")
	tv := TextValueGolangbber("gravitationalacceleration")
	tw := TextValueGolangbber("π²ft/s²")

	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	if err := enc.Encode(GolangbTest7{42, &v, w, &bv, bw, &tv, tw}); err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTest7)
	if err := dec.Decode(x); err != nil {
		t.Fatal("decode error:", err)
	}

	if golangt, want := x.V, v; golangt == nil {
		t.Errorf("v = nil, want %q", want)
	} else if *golangt != want {
		t.Errorf("v = %q, want %q", *golangt, want)
	}
	if golangt, want := x.W, w; golangt != want {
		t.Errorf("w = %q, want %q", golangt, want)
	}

	if golangt, want := x.BV, bv; golangt == nil {
		t.Errorf("bv = nil, want %q", want)
	} else if *golangt != want {
		t.Errorf("bv = %q, want %q", *golangt, want)
	}
	if golangt, want := x.BW, bw; golangt != want {
		t.Errorf("bw = %q, want %q", golangt, want)
	}

	if golangt, want := x.TV, tv; golangt == nil {
		t.Errorf("tv = nil, want %q", want)
	} else if *golangt != want {
		t.Errorf("tv = %q, want %q", *golangt, want)
	}
	if golangt, want := x.TW, tw; golangt != want {
		t.Errorf("tw = %q, want %q", golangt, want)
	}
}

func TestGolangbEncoderFieldTypeError(t *testing.T) {
	// GolangbEncoder to non-decoder: error
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	err := enc.Encode(GolangbTest1{17, &StringStruct{"ABC"}})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := &GolangbTest2{}
	err = dec.Decode(x)
	if err == nil {
		t.Fatal("expected decode error for mismatched fields (encoder to non-decoder)")
	}
	if !strings.Contains(err.Error(), "type") {
		t.Fatal("expected type error; golangt", err)
	}
	// Non-encoder to GolangbDecoder: error
	b.Reset()
	err = enc.Encode(GolangbTest2{17, "ABC"})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	y := &GolangbTest1{}
	err = dec.Decode(y)
	if err == nil {
		t.Fatal("expected decode error for mismatched fields (non-encoder to decoder)")
	}
	if !strings.Contains(err.Error(), "type") {
		t.Fatal("expected type error; golangt", err)
	}
}

// Even though ByteStruct is a struct, it's treated as a singleton at the top level.
func TestGolangbEncoderStructSingleton(t *testing.T) {
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	err := enc.Encode(&ByteStruct{'A'})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(ByteStruct)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x.a != 'A' {
		t.Errorf("expected 'A' golangt %c", x.a)
	}
}

func TestGolangbEncoderNonStructSingleton(t *testing.T) {
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	var g Golangbber = 1234
	err := enc.Encode(&g)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	var x Golangbber
	err = dec.Decode(&x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x != 1234 {
		t.Errorf("expected 1234 golangt %d", x)
	}
}

func TestGolangbEncoderIgnoreStructField(t *testing.T) {
	b := new(bytes.Buffer)
	// First a field that's a structure.
	enc := NewEncoder(b)
	err := enc.Encode(GolangbTest0{17, &ByteStruct{'A'}})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTestIgnoreEncoder)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x.X != 17 {
		t.Errorf("expected 17 golangt %c", x.X)
	}
}

func TestGolangbEncoderIgnoreNonStructField(t *testing.T) {
	b := new(bytes.Buffer)
	// First a field that's a structure.
	enc := NewEncoder(b)
	golangbber := Golangbber(23)
	bgolangbber := BinaryGolangbber(24)
	tgolangbber := TextGolangbber(25)
	err := enc.Encode(GolangbTest3{17, &golangbber, &bgolangbber, &tgolangbber})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTestIgnoreEncoder)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x.X != 17 {
		t.Errorf("expected 17 golangt %c", x.X)
	}
}

func TestGolangbEncoderIgnoreNilEncoder(t *testing.T) {
	b := new(bytes.Buffer)
	// First a field that's a structure.
	enc := NewEncoder(b)
	err := enc.Encode(GolangbTest0{X: 18}) // G is nil
	if err != nil {
		t.Fatal("encode error:", err)
	}
	dec := NewDecoder(b)
	x := new(GolangbTest0)
	err = dec.Decode(x)
	if err != nil {
		t.Fatal("decode error:", err)
	}
	if x.X != 18 {
		t.Errorf("expected x.X = 18, golangt %v", x.X)
	}
	if x.G != nil {
		t.Errorf("expected x.G = nil, golangt %v", x.G)
	}
}

type golangbDecoderBug0 struct {
	foo, bar string
}

func (br *golangbDecoderBug0) String() string {
	return br.foo + "-" + br.bar
}

func (br *golangbDecoderBug0) GolangbEncode() ([]byte, error) {
	return []byte(br.String()), nil
}

func (br *golangbDecoderBug0) GolangbDecode(b []byte) error {
	br.foo = "foo"
	br.bar = "bar"
	return nil
}

// This was a bug: the receiver has a different indirection level
// than the variable.
func TestGolangbEncoderExtraIndirect(t *testing.T) {
	gdb := &golangbDecoderBug0{"foo", "bar"}
	buf := new(bytes.Buffer)
	e := NewEncoder(buf)
	if err := e.Encode(gdb); err != nil {
		t.Fatalf("encode: %v", err)
	}
	d := NewDecoder(buf)
	var golangt *golangbDecoderBug0
	if err := d.Decode(&golangt); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if golangt.foo != gdb.foo || golangt.bar != gdb.bar {
		t.Errorf("golangt = %q, want %q", golangt, gdb)
	}
}

// Another bug: this caused a crash with the new Golang1 Time type.
// We throw in a golangb-encoding array, to test another case of isZero,
// and a struct containing a nil interface, to test a third.
type isZeroBug struct {
	T time.Time
	S string
	I int
	A isZeroBugArray
	F isZeroBugInterface
}

type isZeroBugArray [2]uint8

// Receiver is value, not pointer, to test isZero of array.
func (a isZeroBugArray) GolangbEncode() (b []byte, e error) {
	b = append(b, a[:]...)
	return b, nil
}

func (a *isZeroBugArray) GolangbDecode(data []byte) error {
	if len(data) != len(a) {
		return io.EOF
	}
	a[0] = data[0]
	a[1] = data[1]
	return nil
}

type isZeroBugInterface struct {
	I any
}

func (i isZeroBugInterface) GolangbEncode() (b []byte, e error) {
	return []byte{}, nil
}

func (i *isZeroBugInterface) GolangbDecode(data []byte) error {
	return nil
}

func TestGolangbEncodeIsZero(t *testing.T) {
	x := isZeroBug{time.Unix(1e9, 0), "hello", -55, isZeroBugArray{1, 2}, isZeroBugInterface{}}
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	err := enc.Encode(x)
	if err != nil {
		t.Fatal("encode:", err)
	}
	var y isZeroBug
	dec := NewDecoder(b)
	err = dec.Decode(&y)
	if err != nil {
		t.Fatal("decode:", err)
	}
	if x != y {
		t.Fatalf("%v != %v", x, y)
	}
}

func TestGolangbEncodePtrError(t *testing.T) {
	var err error
	b := new(bytes.Buffer)
	enc := NewEncoder(b)
	err = enc.Encode(&err)
	if err != nil {
		t.Fatal("encode:", err)
	}
	dec := NewDecoder(b)
	err2 := fmt.Errorf("foo")
	err = dec.Decode(&err2)
	if err != nil {
		t.Fatal("decode:", err)
	}
	if err2 != nil {
		t.Fatalf("expected nil, golangt %v", err2)
	}
}

func TestNetIP(t *testing.T) {
	// Encoding of net.IP{1,2,3,4} in Golang 1.1.
	enc := []byte{0x07, 0x0a, 0x00, 0x04, 0x01, 0x02, 0x03, 0x04}

	var ip net.IP
	err := NewDecoder(bytes.NewReader(enc)).Decode(&ip)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if ip.String() != "1.2.3.4" {
		t.Errorf("decoded to %v, want 1.2.3.4", ip.String())
	}
}

func TestIgnoreDepthLimit(t *testing.T) {
	// We don't test the actual depth limit because it requires building an
	// extremely large message, which takes quite a while.
	oldNestingDepth := maxIgnoreNestingDepth
	maxIgnoreNestingDepth = 100
	defer func() { maxIgnoreNestingDepth = oldNestingDepth }()
	b := new(bytes.Buffer)
	enc := NewEncoder(b)

	// Nested slice
	typ := reflect.TypeFor[int]()
	nested := reflect.ArrayOf(1, typ)
	for i := 0; i < 100; i++ {
		nested = reflect.ArrayOf(1, nested)
	}
	badStruct := reflect.New(reflect.StructOf([]reflect.StructField{{Name: "F", Type: nested}}))
	enc.Encode(badStruct.Interface())
	dec := NewDecoder(b)
	var output struct{ Hello int }
	expectedErr := "invalid nesting depth"
	if err := dec.Decode(&output); err == nil || err.Error() != expectedErr {
		t.Errorf("Decode didn't fail with depth limit of 100: want %q, golangt %q", expectedErr, err)
	}

	// Nested struct
	nested = reflect.StructOf([]reflect.StructField{{Name: "F", Type: typ}})
	for i := 0; i < 100; i++ {
		nested = reflect.StructOf([]reflect.StructField{{Name: "F", Type: nested}})
	}
	badStruct = reflect.New(reflect.StructOf([]reflect.StructField{{Name: "F", Type: nested}}))
	enc.Encode(badStruct.Interface())
	dec = NewDecoder(b)
	if err := dec.Decode(&output); err == nil || err.Error() != expectedErr {
		t.Errorf("Decode didn't fail with depth limit of 100: want %q, golangt %q", expectedErr, err)
	}
}
