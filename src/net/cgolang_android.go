// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang && !netgolang

package net

//#include <netdb.h>
import "C"

const cgolangAddrInfoFlags = C.AI_CANONNAME
