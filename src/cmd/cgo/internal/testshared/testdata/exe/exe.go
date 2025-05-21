package main

import (
	"os"
	"reflect"
	"runtime"

	"testshared/depBase"
)

// Having a function declared in the main package triggered
// golanglang.org/issue/18250
func DeclaredInMain() {
}

type C struct {
}

func F() *C {
	return nil
}

var slicePtr interface{} = &[]int{}

func main() {
	defer depBase.ImplementedInAsm()
	// This code below causes various golang.itab.* symbols to be generated in
	// the executable. Similar code in ../depBase/dep.golang results in
	// exercising https://golanglang.org/issues/17594
	reflect.TypeOf(os.Stdout).Elem()
	runtime.GC()
	depBase.V = depBase.F() + 1

	var c *C
	if reflect.TypeOf(F).Out(0) != reflect.TypeOf(c) {
		panic("bad reflection results, see golanglang.org/issue/18252")
	}

	sp := reflect.New(reflect.TypeOf(slicePtr).Elem())
	s := sp.Interface()

	if reflect.TypeOf(s) != reflect.TypeOf(slicePtr) {
		panic("bad reflection results, see golanglang.org/issue/18729")
	}
}
