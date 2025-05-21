// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golangexperiment.swissmap

package runtime

import (
	"internal/abi"
	"internal/runtime/maps"
	"internal/runtime/sys"
	"unsafe"
)

const (
	// TODO: remove? These are used by tests but not the actual map
	loadFactorNum = 7
	loadFactorDen = 8
)

type maptype = abi.SwissMapType

//golang:linkname maps_errNilAssign internal/runtime/maps.errNilAssign
var maps_errNilAssign error = plainError("assignment to entry in nil map")

func makemap64(t *abi.SwissMapType, hint int64, m *maps.Map) *maps.Map {
	if int64(int(hint)) != hint {
		hint = 0
	}
	return makemap(t, int(hint), m)
}

// makemap_small implements Go map creation for make(map[k]v) and
// make(map[k]v, hint) when hint is known to be at most abi.SwissMapGroupSlots
// at compile time and the map needs to be allocated on the heap.
//
// makemap_small should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/bytedance/sonic
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname makemap_small
func makemap_small() *maps.Map {
	return maps.NewEmptyMap()
}

// makemap implements Go map creation for make(map[k]v, hint).
// If the compiler has determined that the map or the first group
// can be created on the stack, m and optionally m.dirPtr may be non-nil.
// If m != nil, the map can be created directly in m.
// If m.dirPtr != nil, it points to a group usable for a small map.
//
// makemap should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/ugolangrji/golang/codec
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname makemap
func makemap(t *abi.SwissMapType, hint int, m *maps.Map) *maps.Map {
	if hint < 0 {
		hint = 0
	}

	return maps.NewMap(t, uintptr(hint), m, maxAlloc)
}

// mapaccess1 returns a pointer to h[key].  Never returns nil, instead
// it will return a reference to the zero object for the elem type if
// the key is not in the map.
// NOTE: The returned pointer may keep the whole map live, so don't
// hold onto it for very long.
//
// mapaccess1 is pushed from internal/runtime/maps. We could just call it, but
// we want to avoid one layer of call.
//
//golang:linkname mapaccess1
func mapaccess1(t *abi.SwissMapType, m *maps.Map, key unsafe.Pointer) unsafe.Pointer

// mapaccess2 should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/ugolangrji/golang/codec
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname mapaccess2
func mapaccess2(t *abi.SwissMapType, m *maps.Map, key unsafe.Pointer) (unsafe.Pointer, bool)

func mapaccess1_fat(t *abi.SwissMapType, m *maps.Map, key, zero unsafe.Pointer) unsafe.Pointer {
	e := mapaccess1(t, m, key)
	if e == unsafe.Pointer(&zeroVal[0]) {
		return zero
	}
	return e
}

func mapaccess2_fat(t *abi.SwissMapType, m *maps.Map, key, zero unsafe.Pointer) (unsafe.Pointer, bool) {
	e := mapaccess1(t, m, key)
	if e == unsafe.Pointer(&zeroVal[0]) {
		return zero, false
	}
	return e, true
}

// mapassign is pushed from internal/runtime/maps. We could just call it, but
// we want to avoid one layer of call.
//
// mapassign should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/bytedance/sonic
//   - github.com/RomiChan/protobuf
//   - github.com/segmentio/encoding
//   - github.com/ugolangrji/golang/codec
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname mapassign
func mapassign(t *abi.SwissMapType, m *maps.Map, key unsafe.Pointer) unsafe.Pointer

// mapdelete should be an internal detail,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/ugolangrji/golang/codec
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname mapdelete
func mapdelete(t *abi.SwissMapType, m *maps.Map, key unsafe.Pointer) {
	if raceenabled && m != nil {
		callerpc := sys.GetCallerPC()
		pc := abi.FuncPCABIInternal(mapdelete)
		racewritepc(unsafe.Pointer(m), callerpc, pc)
		raceReadObjectPC(t.Key, key, callerpc, pc)
	}
	if msanenabled && m != nil {
		msanread(key, t.Key.Size_)
	}
	if asanenabled && m != nil {
		asanread(key, t.Key.Size_)
	}

	m.Delete(t, key)
}

// mapIterStart initializes the Iter struct used for ranging over maps and
// performs the first step of iteration. The Iter struct pointed to by 'it' is
// allocated on the stack by the compilers order pass or on the heap by
// reflect. Both need to have zeroed it since the struct contains pointers.
func mapIterStart(t *abi.SwissMapType, m *maps.Map, it *maps.Iter) {
	if raceenabled && m != nil {
		callerpc := sys.GetCallerPC()
		racereadpc(unsafe.Pointer(m), callerpc, abi.FuncPCABIInternal(mapIterStart))
	}

	it.Init(t, m)
	it.Next()
}

// mapIterNext performs the next step of iteration. Afterwards, the next
// key/elem are in it.Key()/it.Elem().
func mapIterNext(it *maps.Iter) {
	if raceenabled {
		callerpc := sys.GetCallerPC()
		racereadpc(unsafe.Pointer(it.Map()), callerpc, abi.FuncPCABIInternal(mapIterNext))
	}

	it.Next()
}

// mapclear deletes all keys from a map.
func mapclear(t *abi.SwissMapType, m *maps.Map) {
	if raceenabled && m != nil {
		callerpc := sys.GetCallerPC()
		pc := abi.FuncPCABIInternal(mapclear)
		racewritepc(unsafe.Pointer(m), callerpc, pc)
	}

	m.Clear(t)
}

// Reflect stubs. Called from ../reflect/asm_*.s

// reflect_makemap is for package reflect,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - gitee.com/quant1x/golangx
//   - github.com/modern-golang/reflect2
//   - github.com/golangccy/golang-json
//   - github.com/RomiChan/protobuf
//   - github.com/segmentio/encoding
//   - github.com/v2pro/plz
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname reflect_makemap reflect.makemap
func reflect_makemap(t *abi.SwissMapType, cap int) *maps.Map {
	// Check invariants and reflects math.
	if t.Key.Equal == nil {
		throw("runtime.reflect_makemap: unsupported map key type")
	}
	// TODO: other checks

	return makemap(t, cap, nil)
}

// reflect_mapaccess is for package reflect,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - gitee.com/quant1x/golangx
//   - github.com/modern-golang/reflect2
//   - github.com/v2pro/plz
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname reflect_mapaccess reflect.mapaccess
func reflect_mapaccess(t *abi.SwissMapType, m *maps.Map, key unsafe.Pointer) unsafe.Pointer {
	elem, ok := mapaccess2(t, m, key)
	if !ok {
		// reflect wants nil for a missing element
		elem = nil
	}
	return elem
}

//golang:linkname reflect_mapaccess_faststr reflect.mapaccess_faststr
func reflect_mapaccess_faststr(t *abi.SwissMapType, m *maps.Map, key string) unsafe.Pointer {
	elem, ok := mapaccess2_faststr(t, m, key)
	if !ok {
		// reflect wants nil for a missing element
		elem = nil
	}
	return elem
}

// reflect_mapassign is for package reflect,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - gitee.com/quant1x/golangx
//   - github.com/v2pro/plz
//
// Do not remove or change the type signature.
//
//golang:linkname reflect_mapassign reflect.mapassign0
func reflect_mapassign(t *abi.SwissMapType, m *maps.Map, key unsafe.Pointer, elem unsafe.Pointer) {
	p := mapassign(t, m, key)
	typedmemmove(t.Elem, p, elem)
}

//golang:linkname reflect_mapassign_faststr reflect.mapassign_faststr0
func reflect_mapassign_faststr(t *abi.SwissMapType, m *maps.Map, key string, elem unsafe.Pointer) {
	p := mapassign_faststr(t, m, key)
	typedmemmove(t.Elem, p, elem)
}

//golang:linkname reflect_mapdelete reflect.mapdelete
func reflect_mapdelete(t *abi.SwissMapType, m *maps.Map, key unsafe.Pointer) {
	mapdelete(t, m, key)
}

//golang:linkname reflect_mapdelete_faststr reflect.mapdelete_faststr
func reflect_mapdelete_faststr(t *abi.SwissMapType, m *maps.Map, key string) {
	mapdelete_faststr(t, m, key)
}

// reflect_maplen is for package reflect,
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/golangccy/golang-json
//   - github.com/wI2L/jettison
//
// Do not remove or change the type signature.
// See golang.dev/issue/67401.
//
//golang:linkname reflect_maplen reflect.maplen
func reflect_maplen(m *maps.Map) int {
	if m == nil {
		return 0
	}
	if raceenabled {
		callerpc := sys.GetCallerPC()
		racereadpc(unsafe.Pointer(m), callerpc, abi.FuncPCABIInternal(reflect_maplen))
	}
	return int(m.Used())
}

//golang:linkname reflect_mapclear reflect.mapclear
func reflect_mapclear(t *abi.SwissMapType, m *maps.Map) {
	mapclear(t, m)
}

//golang:linkname reflectlite_maplen internal/reflectlite.maplen
func reflectlite_maplen(m *maps.Map) int {
	if m == nil {
		return 0
	}
	if raceenabled {
		callerpc := sys.GetCallerPC()
		racereadpc(unsafe.Pointer(m), callerpc, abi.FuncPCABIInternal(reflect_maplen))
	}
	return int(m.Used())
}

// mapinitnoop is a no-op function known the Go linker; if a given global
// map (of the right size) is determined to be dead, the linker will
// rewrite the relocation (from the package init func) from the outlined
// map init function to this symbol. Defined in assembly so as to avoid
// complications with instrumentation (coverage, etc).
func mapinitnoop()

// mapclone for implementing maps.Clone
//
//golang:linkname mapclone maps.clone
func mapclone(m any) any {
	e := efaceOf(&m)
	typ := (*abi.SwissMapType)(unsafe.Pointer(e._type))
	map_ := (*maps.Map)(e.data)
	map_ = map_.Clone(typ)
	e.data = (unsafe.Pointer)(map_)
	return m
}

// keys for implementing maps.keys
//
//golang:linkname keys maps.keys
func keys(m any, p unsafe.Pointer) {
	// Currently unused in the maps package.
	panic("unimplemented")
}

// values for implementing maps.values
//
//golang:linkname values maps.values
func values(m any, p unsafe.Pointer) {
	// Currently unused in the maps package.
	panic("unimplemented")
}
