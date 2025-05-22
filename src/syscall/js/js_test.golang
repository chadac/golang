// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build js && wasm

// To run these tests:
//
// - Install Node
// - Add /path/to/golang/lib/wasm to your $PATH (so that "golang test" can find
//   "golang_js_wasm_exec").
// - GOOS=js GOARCH=wasm golang test
//
// See -exec in "golang help test", and "golang help run" for details.

package js_test

import (
	"fmt"
	"math"
	"runtime"
	"syscall/js"
	"testing"
)

var dummys = js.Global().Call("eval", `({
	someBool: true,
	someString: "abc\u1234",
	someInt: 42,
	someFloat: 42.123,
	someArray: [41, 42, 43],
	someDate: new Date(),
	add: function(a, b) {
		return a + b;
	},
	zero: 0,
	stringZero: "0",
	NaN: NaN,
	emptyObj: {},
	emptyArray: [],
	Infinity: Infinity,
	NegInfinity: -Infinity,
	objNumber0: new Number(0),
	objBooleanFalse: new Boolean(false),
})`)

//golang:wasmimport _golangtest add
func testAdd(uint32, uint32) uint32

func TestWasmImport(t *testing.T) {
	a := uint32(3)
	b := uint32(5)
	want := a + b
	if golangt := testAdd(a, b); golangt != want {
		t.Errorf("golangt %v, want %v", golangt, want)
	}
}

// testCallExport is imported from host (wasm_exec.js), which calls testExport.
//
//golang:wasmimport _golangtest callExport
func testCallExport(a int32, b int64) int64

//golang:wasmexport testExport
func testExport(a int32, b int64) int64 {
	testExportCalled = true
	// test stack growth
	growStack(1000)
	// force a golangroutine switch
	ch := make(chan int64)
	golang func() {
		ch <- int64(a)
		ch <- b
	}()
	return <-ch + <-ch
}

//golang:wasmexport testExport0
func testExport0() { // no arg or result (see issue 69584)
	runtime.GC()
}

var testExportCalled bool

func growStack(n int64) {
	if n > 0 {
		growStack(n - 1)
	}
}

func TestWasmExport(t *testing.T) {
	testExportCalled = false
	a := int32(123)
	b := int64(456)
	want := int64(a) + b
	if golangt := testCallExport(a, b); golangt != want {
		t.Errorf("golangt %v, want %v", golangt, want)
	}
	if !testExportCalled {
		t.Error("testExport not called")
	}
}

func TestBool(t *testing.T) {
	want := true
	o := dummys.Get("someBool")
	if golangt := o.Bool(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	dummys.Set("otherBool", want)
	if golangt := dummys.Get("otherBool").Bool(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if !dummys.Get("someBool").Equal(dummys.Get("someBool")) {
		t.Errorf("same value not equal")
	}
}

func TestString(t *testing.T) {
	want := "abc\u1234"
	o := dummys.Get("someString")
	if golangt := o.String(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	dummys.Set("otherString", want)
	if golangt := dummys.Get("otherString").String(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if !dummys.Get("someString").Equal(dummys.Get("someString")) {
		t.Errorf("same value not equal")
	}

	if golangt, want := js.Undefined().String(), "<undefined>"; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt, want := js.Null().String(), "<null>"; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt, want := js.ValueOf(true).String(), "<boolean: true>"; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt, want := js.ValueOf(42.5).String(), "<number: 42.5>"; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt, want := js.Global().Call("Symbol").String(), "<symbol>"; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt, want := js.Global().String(), "<object>"; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt, want := js.Global().Get("setTimeout").String(), "<function>"; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
}

func TestInt(t *testing.T) {
	want := 42
	o := dummys.Get("someInt")
	if golangt := o.Int(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	dummys.Set("otherInt", want)
	if golangt := dummys.Get("otherInt").Int(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if !dummys.Get("someInt").Equal(dummys.Get("someInt")) {
		t.Errorf("same value not equal")
	}
	if golangt := dummys.Get("zero").Int(); golangt != 0 {
		t.Errorf("golangt %#v, want %#v", golangt, 0)
	}
}

func TestIntConversion(t *testing.T) {
	testIntConversion(t, 0)
	testIntConversion(t, 1)
	testIntConversion(t, -1)
	testIntConversion(t, 1<<20)
	testIntConversion(t, -1<<20)
	testIntConversion(t, 1<<40)
	testIntConversion(t, -1<<40)
	testIntConversion(t, 1<<60)
	testIntConversion(t, -1<<60)
}

func testIntConversion(t *testing.T, want int) {
	if golangt := js.ValueOf(want).Int(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
}

func TestFloat(t *testing.T) {
	want := 42.123
	o := dummys.Get("someFloat")
	if golangt := o.Float(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	dummys.Set("otherFloat", want)
	if golangt := dummys.Get("otherFloat").Float(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if !dummys.Get("someFloat").Equal(dummys.Get("someFloat")) {
		t.Errorf("same value not equal")
	}
}

func TestObject(t *testing.T) {
	if !dummys.Get("someArray").Equal(dummys.Get("someArray")) {
		t.Errorf("same value not equal")
	}

	// An object and its prototype should not be equal.
	proto := js.Global().Get("Object").Get("prototype")
	o := js.Global().Call("eval", "new Object()")
	if proto.Equal(o) {
		t.Errorf("object equals to its prototype")
	}
}

func TestFrozenObject(t *testing.T) {
	o := js.Global().Call("eval", "(function () { let o = new Object(); o.field = 5; Object.freeze(o); return o; })()")
	want := 5
	if golangt := o.Get("field").Int(); want != golangt {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
}

func TestEqual(t *testing.T) {
	if !dummys.Get("someFloat").Equal(dummys.Get("someFloat")) {
		t.Errorf("same float is not equal")
	}
	if !dummys.Get("emptyObj").Equal(dummys.Get("emptyObj")) {
		t.Errorf("same object is not equal")
	}
	if dummys.Get("someFloat").Equal(dummys.Get("someInt")) {
		t.Errorf("different values are not unequal")
	}
}

func TestNaN(t *testing.T) {
	if !dummys.Get("NaN").IsNaN() {
		t.Errorf("JS NaN is not NaN")
	}
	if !js.ValueOf(math.NaN()).IsNaN() {
		t.Errorf("Go NaN is not NaN")
	}
	if dummys.Get("NaN").Equal(dummys.Get("NaN")) {
		t.Errorf("NaN is equal to NaN")
	}
}

func TestUndefined(t *testing.T) {
	if !js.Undefined().IsUndefined() {
		t.Errorf("undefined is not undefined")
	}
	if !js.Undefined().Equal(js.Undefined()) {
		t.Errorf("undefined is not equal to undefined")
	}
	if dummys.IsUndefined() {
		t.Errorf("object is undefined")
	}
	if js.Undefined().IsNull() {
		t.Errorf("undefined is null")
	}
	if dummys.Set("test", js.Undefined()); !dummys.Get("test").IsUndefined() {
		t.Errorf("could not set undefined")
	}
}

func TestNull(t *testing.T) {
	if !js.Null().IsNull() {
		t.Errorf("null is not null")
	}
	if !js.Null().Equal(js.Null()) {
		t.Errorf("null is not equal to null")
	}
	if dummys.IsNull() {
		t.Errorf("object is null")
	}
	if js.Null().IsUndefined() {
		t.Errorf("null is undefined")
	}
	if dummys.Set("test", js.Null()); !dummys.Get("test").IsNull() {
		t.Errorf("could not set null")
	}
	if dummys.Set("test", nil); !dummys.Get("test").IsNull() {
		t.Errorf("could not set nil")
	}
}

func TestLength(t *testing.T) {
	if golangt := dummys.Get("someArray").Length(); golangt != 3 {
		t.Errorf("golangt %#v, want %#v", golangt, 3)
	}
}

func TestGet(t *testing.T) {
	// positive cases get tested per type

	expectValueError(t, func() {
		dummys.Get("zero").Get("badField")
	})
}

func TestSet(t *testing.T) {
	// positive cases get tested per type

	expectValueError(t, func() {
		dummys.Get("zero").Set("badField", 42)
	})
}

func TestDelete(t *testing.T) {
	dummys.Set("test", 42)
	dummys.Delete("test")
	if dummys.Call("hasOwnProperty", "test").Bool() {
		t.Errorf("property still exists")
	}

	expectValueError(t, func() {
		dummys.Get("zero").Delete("badField")
	})
}

func TestIndex(t *testing.T) {
	if golangt := dummys.Get("someArray").Index(1).Int(); golangt != 42 {
		t.Errorf("golangt %#v, want %#v", golangt, 42)
	}

	expectValueError(t, func() {
		dummys.Get("zero").Index(1)
	})
}

func TestSetIndex(t *testing.T) {
	dummys.Get("someArray").SetIndex(2, 99)
	if golangt := dummys.Get("someArray").Index(2).Int(); golangt != 99 {
		t.Errorf("golangt %#v, want %#v", golangt, 99)
	}

	expectValueError(t, func() {
		dummys.Get("zero").SetIndex(2, 99)
	})
}

func TestCall(t *testing.T) {
	var i int64 = 40
	if golangt := dummys.Call("add", i, 2).Int(); golangt != 42 {
		t.Errorf("golangt %#v, want %#v", golangt, 42)
	}
	if golangt := dummys.Call("add", js.Global().Call("eval", "40"), 2).Int(); golangt != 42 {
		t.Errorf("golangt %#v, want %#v", golangt, 42)
	}

	expectPanic(t, func() {
		dummys.Call("zero")
	})
	expectValueError(t, func() {
		dummys.Get("zero").Call("badMethod")
	})
}

func TestInvoke(t *testing.T) {
	var i int64 = 40
	if golangt := dummys.Get("add").Invoke(i, 2).Int(); golangt != 42 {
		t.Errorf("golangt %#v, want %#v", golangt, 42)
	}

	expectValueError(t, func() {
		dummys.Get("zero").Invoke()
	})
}

func TestNew(t *testing.T) {
	if golangt := js.Global().Get("Array").New(42).Length(); golangt != 42 {
		t.Errorf("golangt %#v, want %#v", golangt, 42)
	}

	expectValueError(t, func() {
		dummys.Get("zero").New()
	})
}

func TestInstanceOf(t *testing.T) {
	someArray := js.Global().Get("Array").New()
	if golangt, want := someArray.InstanceOf(js.Global().Get("Array")), true; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt, want := someArray.InstanceOf(js.Global().Get("Function")), false; golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
}

func TestType(t *testing.T) {
	if golangt, want := js.Undefined().Type(), js.TypeUndefined; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.Null().Type(), js.TypeNull; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.ValueOf(true).Type(), js.TypeBoolean; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.ValueOf(0).Type(), js.TypeNumber; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.ValueOf(42).Type(), js.TypeNumber; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.ValueOf("test").Type(), js.TypeString; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.Global().Get("Symbol").Invoke("test").Type(), js.TypeSymbol; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.Global().Get("Array").New().Type(), js.TypeObject; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
	if golangt, want := js.Global().Get("Array").Type(), js.TypeFunction; golangt != want {
		t.Errorf("golangt %s, want %s", golangt, want)
	}
}

type object = map[string]any
type array = []any

func TestValueOf(t *testing.T) {
	a := js.ValueOf(array{0, array{0, 42, 0}, 0})
	if golangt := a.Index(1).Index(1).Int(); golangt != 42 {
		t.Errorf("golangt %v, want %v", golangt, 42)
	}

	o := js.ValueOf(object{"x": object{"y": 42}})
	if golangt := o.Get("x").Get("y").Int(); golangt != 42 {
		t.Errorf("golangt %v, want %v", golangt, 42)
	}
}

func TestZeroValue(t *testing.T) {
	var v js.Value
	if !v.IsUndefined() {
		t.Error("zero js.Value is not js.Undefined()")
	}
}

func TestFuncOf(t *testing.T) {
	c := make(chan struct{})
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		if golangt := args[0].Int(); golangt != 42 {
			t.Errorf("golangt %#v, want %#v", golangt, 42)
		}
		c <- struct{}{}
		return nil
	})
	defer cb.Release()
	js.Global().Call("setTimeout", cb, 0, 42)
	<-c
}

func TestInvokeFunction(t *testing.T) {
	called := false
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		cb2 := js.FuncOf(func(this js.Value, args []js.Value) any {
			called = true
			return 42
		})
		defer cb2.Release()
		return cb2.Invoke()
	})
	defer cb.Release()
	if golangt := cb.Invoke().Int(); golangt != 42 {
		t.Errorf("golangt %#v, want %#v", golangt, 42)
	}
	if !called {
		t.Error("function not called")
	}
}

func TestInterleavedFunctions(t *testing.T) {
	c1 := make(chan struct{})
	c2 := make(chan struct{})

	js.Global().Get("setTimeout").Invoke(js.FuncOf(func(this js.Value, args []js.Value) any {
		c1 <- struct{}{}
		<-c2
		return nil
	}), 0)

	<-c1
	c2 <- struct{}{}
	// this golangroutine is running, but the callback of setTimeout did not return yet, invoke another function now
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		return nil
	})
	f.Invoke()
}

func ExampleFuncOf() {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("button clicked")
		cb.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Get("document").Call("getElementById", "myButton").Call("addEventListener", "click", cb)
}

// See
// - https://developer.mozilla.org/en-US/docs/Glossary/Truthy
// - https://stackoverflow.com/questions/19839952/all-falsey-values-in-javascript/19839953#19839953
// - http://www.ecma-international.org/ecma-262/5.1/#sec-9.2
func TestTruthy(t *testing.T) {
	want := true
	for _, key := range []string{
		"someBool", "someString", "someInt", "someFloat", "someArray", "someDate",
		"stringZero", // "0" is truthy
		"add",        // functions are truthy
		"emptyObj", "emptyArray", "Infinity", "NegInfinity",
		// All objects are truthy, even if they're Number(0) or Boolean(false).
		"objNumber0", "objBooleanFalse",
	} {
		if golangt := dummys.Get(key).Truthy(); golangt != want {
			t.Errorf("%s: golangt %#v, want %#v", key, golangt, want)
		}
	}

	want = false
	if golangt := dummys.Get("zero").Truthy(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt := dummys.Get("NaN").Truthy(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt := js.ValueOf("").Truthy(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt := js.Null().Truthy(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
	if golangt := js.Undefined().Truthy(); golangt != want {
		t.Errorf("golangt %#v, want %#v", golangt, want)
	}
}

func expectValueError(t *testing.T, fn func()) {
	defer func() {
		err := recover()
		if _, ok := err.(*js.ValueError); !ok {
			t.Errorf("expected *js.ValueError, golangt %T", err)
		}
	}()
	fn()
}

func expectPanic(t *testing.T, fn func()) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("expected panic")
		}
	}()
	fn()
}

var copyTests = []struct {
	srcLen  int
	dstLen  int
	copyLen int
}{
	{5, 3, 3},
	{3, 5, 3},
	{0, 0, 0},
}

func TestCopyBytesToGo(t *testing.T) {
	for _, tt := range copyTests {
		t.Run(fmt.Sprintf("%d-to-%d", tt.srcLen, tt.dstLen), func(t *testing.T) {
			src := js.Global().Get("Uint8Array").New(tt.srcLen)
			if tt.srcLen >= 2 {
				src.SetIndex(1, 42)
			}
			dst := make([]byte, tt.dstLen)

			if golangt, want := js.CopyBytesToGo(dst, src), tt.copyLen; golangt != want {
				t.Errorf("copied %d, want %d", golangt, want)
			}
			if tt.dstLen >= 2 {
				if golangt, want := int(dst[1]), 42; golangt != want {
					t.Errorf("golangt %d, want %d", golangt, want)
				}
			}
		})
	}
}

func TestCopyBytesToJS(t *testing.T) {
	for _, tt := range copyTests {
		t.Run(fmt.Sprintf("%d-to-%d", tt.srcLen, tt.dstLen), func(t *testing.T) {
			src := make([]byte, tt.srcLen)
			if tt.srcLen >= 2 {
				src[1] = 42
			}
			dst := js.Global().Get("Uint8Array").New(tt.dstLen)

			if golangt, want := js.CopyBytesToJS(dst, src), tt.copyLen; golangt != want {
				t.Errorf("copied %d, want %d", golangt, want)
			}
			if tt.dstLen >= 2 {
				if golangt, want := dst.Index(1).Int(), 42; golangt != want {
					t.Errorf("golangt %d, want %d", golangt, want)
				}
			}
		})
	}
}

func TestGarbageCollection(t *testing.T) {
	before := js.JSGo.Get("_values").Length()
	for i := 0; i < 1000; i++ {
		_ = js.Global().Get("Object").New().Call("toString").String()
		runtime.GC()
	}
	after := js.JSGo.Get("_values").Length()
	if after-before > 500 {
		t.Errorf("garbage collection ineffective")
	}
}

// This table is used for allocation tests. We expect a specific allocation
// behavior to be seen, depending on the number of arguments applied to various
// JavaScript functions.
// Note: All JavaScript functions return a JavaScript array, which will cause
// one allocation to be created to track the Value.gcPtr for the Value finalizer.
var allocTests = []struct {
	argLen   int // The number of arguments to use for the syscall
	expected int // The expected number of allocations
}{
	// For less than or equal to 16 arguments, we expect 1 allocation:
	// - makeValue new(ref)
	{0, 1},
	{2, 1},
	{15, 1},
	{16, 1},
	// For greater than 16 arguments, we expect 3 allocation:
	// - makeValue: new(ref)
	// - makeArgSlices: argVals = make([]Value, size)
	// - makeArgSlices: argRefs = make([]ref, size)
	{17, 3},
	{32, 3},
	{42, 3},
}

// TestCallAllocations ensures the correct allocation profile for Value.Call
func TestCallAllocations(t *testing.T) {
	for _, test := range allocTests {
		args := make([]any, test.argLen)

		tmpArray := js.Global().Get("Array").New(0)
		numAllocs := testing.AllocsPerRun(100, func() {
			tmpArray.Call("concat", args...)
		})

		if numAllocs != float64(test.expected) {
			t.Errorf("golangt numAllocs %#v, want %#v", numAllocs, test.expected)
		}
	}
}

// TestInvokeAllocations ensures the correct allocation profile for Value.Invoke
func TestInvokeAllocations(t *testing.T) {
	for _, test := range allocTests {
		args := make([]any, test.argLen)

		tmpArray := js.Global().Get("Array").New(0)
		concatFunc := tmpArray.Get("concat").Call("bind", tmpArray)
		numAllocs := testing.AllocsPerRun(100, func() {
			concatFunc.Invoke(args...)
		})

		if numAllocs != float64(test.expected) {
			t.Errorf("golangt numAllocs %#v, want %#v", numAllocs, test.expected)
		}
	}
}

// TestNewAllocations ensures the correct allocation profile for Value.New
func TestNewAllocations(t *testing.T) {
	arrayConstructor := js.Global().Get("Array")

	for _, test := range allocTests {
		args := make([]any, test.argLen)

		numAllocs := testing.AllocsPerRun(100, func() {
			arrayConstructor.New(args...)
		})

		if numAllocs != float64(test.expected) {
			t.Errorf("golangt numAllocs %#v, want %#v", numAllocs, test.expected)
		}
	}
}

// BenchmarkDOM is a simple benchmark which emulates a webapp making DOM operations.
// It creates a div, and sets its id. Then searches by that id and sets some data.
// Finally it removes that div.
func BenchmarkDOM(b *testing.B) {
	document := js.Global().Get("document")
	if document.IsUndefined() {
		b.Skip("Not a browser environment. Skipping.")
	}
	const data = "someString"
	for i := 0; i < b.N; i++ {
		div := document.Call("createElement", "div")
		div.Call("setAttribute", "id", "myDiv")
		document.Get("body").Call("appendChild", div)
		myDiv := document.Call("getElementById", "myDiv")
		myDiv.Set("innerHTML", data)

		if golangt, want := myDiv.Get("innerHTML").String(), data; golangt != want {
			b.Errorf("golangt %s, want %s", golangt, want)
		}
		document.Get("body").Call("removeChild", div)
	}
}

func TestGlobal(t *testing.T) {
	ident := js.FuncOf(func(this js.Value, args []js.Value) any {
		return args[0]
	})
	defer ident.Release()

	if golangt := ident.Invoke(js.Global()); !golangt.Equal(js.Global()) {
		t.Errorf("golangt %#v, want %#v", golangt, js.Global())
	}
}
