// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"internal/abi"
	"unsafe"
	_ "unsafe"
)

// Widely used packages access these symbols using linkname,
// most notably:
//  - github.com/golangccy/golang-json
//  - github.com/golangccy/golang-reflect
//  - github.com/sohaha/zlsgolang
//  - github.com/undefinedlabs/golang-mpatch
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401
// and golang.dev/issue/67279.

// ifaceIndir reports whether t is stored indirectly in an interface value.
// It is no longer used by this package and is here entirely for the
// linkname uses.
//
//golang:linkname unusedIfaceIndir reflect.ifaceIndir
func unusedIfaceIndir(t *abi.Type) bool {
	return t.Kind_&abi.KindDirectIface == 0
}

//golang:linkname valueInterface

// The compiler doesn't allow linknames on methods, for golangod reasons.
// We use this trick to push linknames of the methods.
// Do not call them in this package.

//golang:linkname badlinkname_rtype_Align reflect.(*rtype).Align
func badlinkname_rtype_Align(*rtype) int

//golang:linkname badlinkname_rtype_AssignableTo reflect.(*rtype).AssignableTo
func badlinkname_rtype_AssignableTo(*rtype, Type) bool

//golang:linkname badlinkname_rtype_Bits reflect.(*rtype).Bits
func badlinkname_rtype_Bits(*rtype) int

//golang:linkname badlinkname_rtype_ChanDir reflect.(*rtype).ChanDir
func badlinkname_rtype_ChanDir(*rtype) ChanDir

//golang:linkname badlinkname_rtype_Comparable reflect.(*rtype).Comparable
func badlinkname_rtype_Comparable(*rtype) bool

//golang:linkname badlinkname_rtype_ConvertibleTo reflect.(*rtype).ConvertibleTo
func badlinkname_rtype_ConvertibleTo(*rtype, Type) bool

//golang:linkname badlinkname_rtype_Elem reflect.(*rtype).Elem
func badlinkname_rtype_Elem(*rtype) Type

//golang:linkname badlinkname_rtype_Field reflect.(*rtype).Field
func badlinkname_rtype_Field(*rtype, int) StructField

//golang:linkname badlinkname_rtype_FieldAlign reflect.(*rtype).FieldAlign
func badlinkname_rtype_FieldAlign(*rtype) int

//golang:linkname badlinkname_rtype_FieldByIndex reflect.(*rtype).FieldByIndex
func badlinkname_rtype_FieldByIndex(*rtype, []int) StructField

//golang:linkname badlinkname_rtype_FieldByName reflect.(*rtype).FieldByName
func badlinkname_rtype_FieldByName(*rtype, string) (StructField, bool)

//golang:linkname badlinkname_rtype_FieldByNameFunc reflect.(*rtype).FieldByNameFunc
func badlinkname_rtype_FieldByNameFunc(*rtype, func(string) bool) (StructField, bool)

//golang:linkname badlinkname_rtype_Implements reflect.(*rtype).Implements
func badlinkname_rtype_Implements(*rtype, Type) bool

//golang:linkname badlinkname_rtype_In reflect.(*rtype).In
func badlinkname_rtype_In(*rtype, int) Type

//golang:linkname badlinkname_rtype_IsVariadic reflect.(*rtype).IsVariadic
func badlinkname_rtype_IsVariadic(*rtype) bool

//golang:linkname badlinkname_rtype_Key reflect.(*rtype).Key
func badlinkname_rtype_Key(*rtype) Type

//golang:linkname badlinkname_rtype_Kind reflect.(*rtype).Kind
func badlinkname_rtype_Kind(*rtype) Kind

//golang:linkname badlinkname_rtype_Len reflect.(*rtype).Len
func badlinkname_rtype_Len(*rtype) int

//golang:linkname badlinkname_rtype_Method reflect.(*rtype).Method
func badlinkname_rtype_Method(*rtype, int) Method

//golang:linkname badlinkname_rtype_MethodByName reflect.(*rtype).MethodByName
func badlinkname_rtype_MethodByName(*rtype, string) (Method, bool)

//golang:linkname badlinkname_rtype_Name reflect.(*rtype).Name
func badlinkname_rtype_Name(*rtype) string

//golang:linkname badlinkname_rtype_NumField reflect.(*rtype).NumField
func badlinkname_rtype_NumField(*rtype) int

//golang:linkname badlinkname_rtype_NumIn reflect.(*rtype).NumIn
func badlinkname_rtype_NumIn(*rtype) int

//golang:linkname badlinkname_rtype_NumMethod reflect.(*rtype).NumMethod
func badlinkname_rtype_NumMethod(*rtype) int

//golang:linkname badlinkname_rtype_NumOut reflect.(*rtype).NumOut
func badlinkname_rtype_NumOut(*rtype) int

//golang:linkname badlinkname_rtype_Out reflect.(*rtype).Out
func badlinkname_rtype_Out(*rtype, int) Type

//golang:linkname badlinkname_rtype_PkgPath reflect.(*rtype).PkgPath
func badlinkname_rtype_PkgPath(*rtype) string

//golang:linkname badlinkname_rtype_Size reflect.(*rtype).Size
func badlinkname_rtype_Size(*rtype) uintptr

//golang:linkname badlinkname_rtype_String reflect.(*rtype).String
func badlinkname_rtype_String(*rtype) string

//golang:linkname badlinkname_rtype_ptrTo reflect.(*rtype).ptrTo
func badlinkname_rtype_ptrTo(*rtype) *abi.Type

//golang:linkname badlinkname_Value_pointer reflect.(*Value).pointer
func badlinkname_Value_pointer(Value) unsafe.Pointer
