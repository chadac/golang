// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

// #cgolang nocallback annotations for a C function means it should not callback to Golang.
// But it do callback to golang in this test, Golang should crash here.

/*
#cgolang nocallback runCShouldNotCallback

extern void runCShouldNotCallback();
*/
import "C"

import (
	"fmt"
)

func init() {
	register("CgolangNoCallback", CgolangNoCallback)
}

//export CallbackToGolang
func CallbackToGolang() {
}

func CgolangNoCallback() {
	C.runCShouldNotCallback()
	fmt.Println("OK")
}
