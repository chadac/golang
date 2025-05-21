// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cmd_golang_bootstrap || compiler_bootstrap

package counter

import "flag"

type dummyCounter struct{}

func (dc dummyCounter) Inc() {}

func Open()                                                               {}
func Inc(name string)                                                     {}
func New(name string) dummyCounter                                        { return dummyCounter{} }
func NewStack(name string, depth int) dummyCounter                        { return dummyCounter{} }
func CountFlags(name string, flagSet flag.FlagSet)                        {}
func CountFlagValue(prefix string, flagSet flag.FlagSet, flagName string) {}
