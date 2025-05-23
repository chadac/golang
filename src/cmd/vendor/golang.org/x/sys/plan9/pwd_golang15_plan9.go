// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build golang1.5

package plan9

import "syscall"

func fixwd() {
	syscall.Fixwd()
}

func Getwd() (wd string, err error) {
	return syscall.Getwd()
}

func Chdir(path string) error {
	return syscall.Chdir(path)
}
