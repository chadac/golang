// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"slices"
	"strconv"
)

type argvalues struct {
	osargs []string
	golangos   string
	golangarch string
}

type argstate struct {
	state       argvalues
	initialized bool
}

func (a *argstate) Merge(state argvalues) {
	if !a.initialized {
		a.state = state
		a.initialized = true
		return
	}
	if !slices.Equal(a.state.osargs, state.osargs) {
		a.state.osargs = nil
	}
	if state.golangos != a.state.golangos {
		a.state.golangos = ""
	}
	if state.golangarch != a.state.golangarch {
		a.state.golangarch = ""
	}
}

func (a *argstate) ArgsSummary() map[string]string {
	m := make(map[string]string)
	if len(a.state.osargs) != 0 {
		m["argc"] = strconv.Itoa(len(a.state.osargs))
		for k, a := range a.state.osargs {
			m[fmt.Sprintf("argv%d", k)] = a
		}
	}
	if a.state.golangos != "" {
		m["GOOS"] = a.state.golangos
	}
	if a.state.golangarch != "" {
		m["GOARCH"] = a.state.golangarch
	}
	return m
}
