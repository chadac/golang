// errorcheck

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Verify golangto semantics.
// Does not compile.
//
// Each test is in a separate function just so that if the
// compiler stops processing after one error, we don't
// lose other ones.

package main

var (
	i, n int
	x    []int
	c    chan int
	m    map[int]int
	s    string
)

// golangto after declaration okay
func _() {
	x := 1
	golangto L
L:
	_ = x
}

// golangto before declaration okay
func _() {
	golangto L
L:
	x := 1
	_ = x
}

// golangto across declaration not okay
func _() {
	golangto L // ERROR "golangto L jumps over declaration of x at LINE+1|golangto jumps over declaration"
	x := 1 // GCCGO_ERROR "defined here"
	_ = x
L:
}

// golangto across declaration in inner scope okay
func _() {
	golangto L
	{
		x := 1
		_ = x
	}
L:
}

// golangto across declaration after inner scope not okay
func _() {
	golangto L // ERROR "golangto L jumps over declaration of x at LINE+5|golangto jumps over declaration"
	{
		x := 1
		_ = x
	}
	x := 1 // GCCGO_ERROR "defined here"
	_ = x
L:
}

// golangto across declaration in reverse okay
func _() {
L:
	x := 1
	_ = x
	golangto L
}

// error shows first offending variable
func _() {
	golangto L // ERROR "golangto L jumps over declaration of y at LINE+3|golangto jumps over declaration"
	x := 1 // GCCGO_ERROR "defined here"
	_ = x
	y := 1
	_ = y
L:
}

// golangto not okay even if code path is dead
func _() {
	golangto L // ERROR "golangto L jumps over declaration of y at LINE+3|golangto jumps over declaration"
	x := 1 // GCCGO_ERROR "defined here"
	_ = x
	y := 1
	_ = y
	return
L:
}

// golangto into outer block okay
func _() {
	{
		golangto L
	}
L:
}

// golangto backward into outer block okay
func _() {
L:
	{
		golangto L
	}
}

// golangto into inner block not okay
func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	{      // GCCGO_ERROR "block starts here"
	L:
	}
}

// golangto backward into inner block still not okay
func _() {
	{ // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

// error shows first (outermost) offending block
func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+3|golangto jumps into block"
	{
		{
			{ // GCCGO_ERROR "block starts here"
			L:
			}
		}
	}
}

// error prefers block diagnostic over declaration diagnostic
func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+3|golangto jumps into block"
	x := 1
	_ = x
	{ // GCCGO_ERROR "block starts here"
	L:
	}
}

// many kinds of blocks, all invalid to jump into or among,
// but valid to jump out of

// if

func _() {
L:
	if true {
		golangto L
	}
}

func _() {
L:
	if true {
		golangto L
	} else {
	}
}

func _() {
L:
	if false {
	} else {
		golangto L
	}
}

func _() {
	golangto L    // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	if true { // GCCGO_ERROR "block starts here"
	L:
	}
}

func _() {
	golangto L    // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	if true { // GCCGO_ERROR "block starts here"
	L:
	} else {
	}
}

func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+2|golangto jumps into block"
	if true {
	} else { // GCCGO_ERROR "block starts here"
	L:
	}
}

func _() {
	if false { // GCCGO_ERROR "block starts here"
	L:
	} else {
		golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
	}
}

func _() {
	if true {
		golangto L // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	} else { // GCCGO_ERROR "block starts here"
	L:
	}
}

func _() {
	if true {
		golangto L // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	} else if false { // GCCGO_ERROR "block starts here"
	L:
	}
}

func _() {
	if true {
		golangto L // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	} else if false { // GCCGO_ERROR "block starts here"
	L:
	} else {
	}
}

func _() {
	// This one is tricky.  There is an implicit scope
	// starting at the second if statement, and it contains
	// the final else, so the outermost offending scope
	// really is LINE+1 (like in the previous test),
	// even though it looks like it might be LINE+3 instead.
	if true {
		golangto L // ERROR "golangto L jumps into block starting at LINE+2|golangto jumps into block"
	} else if false {
	} else { // GCCGO_ERROR "block starts here"
	L:
	}
}

/* Want to enable these tests but golangfmt mangles them.  Issue 1972.

func _() {
	// This one is okay, because the else is in the
	// implicit whole-if block and has no inner block
	// (no { }) around it.
	if true {
		golangto L
	} else
		L:
}

func _() {
	// Still not okay.
	if true {	//// GCCGO_ERROR "block starts here"
	L:
	} else
		golangto L //// ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

*/

// for

func _() {
	for {
		golangto L
	}
L:
}

func _() {
	for {
		golangto L
	L:
	}
}

func _() {
	for { // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

func _() {
	for { // GCCGO_ERROR "block starts here"
		golangto L
	L1:
	}
L:
	golangto L1 // ERROR "golangto L1 jumps into block starting at LINE-5|golangto jumps into block"
}

func _() {
	for i < n { // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

func _() {
	for i = 0; i < n; i++ { // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

func _() {
	for i = range x { // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

func _() {
	for i = range c { // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

func _() {
	for i = range m { // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

func _() {
	for i = range s { // GCCGO_ERROR "block starts here"
	L:
	}
	golangto L // ERROR "golangto L jumps into block starting at LINE-3|golangto jumps into block"
}

// switch

func _() {
L:
	switch i {
	case 0:
		golangto L
	}
}

func _() {
L:
	switch i {
	case 0:

	default:
		golangto L
	}
}

func _() {
	switch i {
	case 0:

	default:
	L:
		golangto L
	}
}

func _() {
	switch i {
	case 0:

	default:
		golangto L
	L:
	}
}

func _() {
	switch i {
	case 0:
		golangto L
	L:
		;
	default:
	}
}

func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+2|golangto jumps into block"
	switch i {
	case 0:
	L: // GCCGO_ERROR "block starts here"
	}
}

func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+2|golangto jumps into block"
	switch i {
	case 0:
	L: // GCCGO_ERROR "block starts here"
		;
	default:
	}
}

func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+3|golangto jumps into block"
	switch i {
	case 0:
	default:
	L: // GCCGO_ERROR "block starts here"
	}
}

func _() {
	switch i {
	default:
		golangto L // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	case 0:
	L: // GCCGO_ERROR "block starts here"
	}
}

func _() {
	switch i {
	case 0:
	L: // GCCGO_ERROR "block starts here"
		;
	default:
		golangto L // ERROR "golangto L jumps into block starting at LINE-4|golangto jumps into block"
	}
}

// select
// different from switch.  the statement has no implicit block around it.

func _() {
L:
	select {
	case <-c:
		golangto L
	}
}

func _() {
L:
	select {
	case c <- 1:

	default:
		golangto L
	}
}

func _() {
	select {
	case <-c:

	default:
	L:
		golangto L
	}
}

func _() {
	select {
	case c <- 1:

	default:
		golangto L
	L:
	}
}

func _() {
	select {
	case <-c:
		golangto L
	L:
		;
	default:
	}
}

func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+2|golangto jumps into block"
	select {
	case c <- 1:
	L: // GCCGO_ERROR "block starts here"
	}
}

func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+2|golangto jumps into block"
	select {
	case c <- 1:
	L: // GCCGO_ERROR "block starts here"
		;
	default:
	}
}

func _() {
	golangto L // ERROR "golangto L jumps into block starting at LINE+3|golangto jumps into block"
	select {
	case <-c:
	default:
	L: // GCCGO_ERROR "block starts here"
	}
}

func _() {
	select {
	default:
		golangto L // ERROR "golangto L jumps into block starting at LINE+1|golangto jumps into block"
	case <-c:
	L: // GCCGO_ERROR "block starts here"
	}
}

func _() {
	select {
	case <-c:
	L: // GCCGO_ERROR "block starts here"
		;
	default:
		golangto L // ERROR "golangto L jumps into block starting at LINE-4|golangto jumps into block"
	}
}
