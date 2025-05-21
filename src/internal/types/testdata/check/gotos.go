// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This file is a modified copy of $GOROOT/test/golangto.golang.

package golangtos

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
	golangto L /* ERROR "golangto L jumps over variable declaration at line 36" */
	x := 1
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
	golangto L /* ERROR "golangto L jumps over variable declaration at line 58" */
	{
		x := 1
		_ = x
	}
	x := 1
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

func _() {
L: L1:
	x := 1
	_ = x
	golangto L
	golangto L1
}

// error shows first offending variable
func _() {
	golangto L /* ERROR "golangto L jumps over variable declaration at line 84" */
	x := 1
	_ = x
	y := 1
	_ = y
L:
}

// golangto not okay even if code path is dead
func _() {
	golangto L /* ERROR "golangto L jumps over variable declaration" */
	x := 1
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

func _() {
	{
		golangto L
		golangto L1
	}
L: L1:
}

// golangto backward into outer block okay
func _() {
L:
	{
		golangto L
	}
}

func _() {
L: L1:
	{
		golangto L
		golangto L1
	}
}

// golangto into inner block not okay
func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	{
	L:
	}
}

func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	golangto L1 /* ERROR "golangto L1 jumps into block" */
	{
	L: L1:
	}
}

// golangto backward into inner block still not okay
func _() {
	{
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
}

func _() {
	{
	L: L1:
	}
	golangto L /* ERROR "golangto L jumps into block" */
	golangto L1 /* ERROR "golangto L1 jumps into block" */
}

// error shows first (outermost) offending block
func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	{
		{
			{
			L:
			}
		}
	}
}

// error prefers block diagnostic over declaration diagnostic
func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	x := 1
	_ = x
	{
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
	golangto L /* ERROR "golangto L jumps into block" */
	if true {
	L:
	}
}

func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	if true {
	L:
	} else {
	}
}

func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	if true {
	} else {
	L:
	}
}

func _() {
	if false {
	L:
	} else {
		golangto L /* ERROR "golangto L jumps into block" */
	}
}

func _() {
	if true {
		golangto L /* ERROR "golangto L jumps into block" */
	} else {
	L:
	}
}

func _() {
	if true {
		golangto L /* ERROR "golangto L jumps into block" */
	} else if false {
	L:
	}
}

func _() {
	if true {
		golangto L /* ERROR "golangto L jumps into block" */
	} else if false {
	L:
	} else {
	}
}

func _() {
	if true {
		golangto L /* ERROR "golangto L jumps into block" */
	} else if false {
	} else {
	L:
	}
}

func _() {
	if true {
		golangto L /* ERROR "golangto L jumps into block" */
	} else {
		L:
	}
}

func _() {
	if true {
		L:
	} else {
		golangto L /* ERROR "golangto L jumps into block" */
	}
}

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
	for {
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
}

func _() {
	for {
		golangto L
	L1:
	}
L:
	golangto L1 /* ERROR "golangto L1 jumps into block" */
}

func _() {
	for i < n {
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
}

func _() {
	for i = 0; i < n; i++ {
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
}

func _() {
	for i = range x {
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
}

func _() {
	for i = range c {
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
}

func _() {
	for i = range m {
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
}

func _() {
	for i = range s {
	L:
	}
	golangto L /* ERROR "golangto L jumps into block" */
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
	golangto L /* ERROR "golangto L jumps into block" */
	switch i {
	case 0:
	L:
	}
}

func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	switch i {
	case 0:
	L:
		;
	default:
	}
}

func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	switch i {
	case 0:
	default:
	L:
	}
}

func _() {
	switch i {
	default:
		golangto L /* ERROR "golangto L jumps into block" */
	case 0:
	L:
	}
}

func _() {
	switch i {
	case 0:
	L:
		;
	default:
		golangto L /* ERROR "golangto L jumps into block" */
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
	golangto L /* ERROR "golangto L jumps into block" */
	select {
	case c <- 1:
	L:
	}
}

func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	select {
	case c <- 1:
	L:
		;
	default:
	}
}

func _() {
	golangto L /* ERROR "golangto L jumps into block" */
	select {
	case <-c:
	default:
	L:
	}
}

func _() {
	select {
	default:
		golangto L /* ERROR "golangto L jumps into block" */
	case <-c:
	L:
	}
}

func _() {
	select {
	case <-c:
	L:
		;
	default:
		golangto L /* ERROR "golangto L jumps into block" */
	}
}
