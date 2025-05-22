// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang && !netgolang && (dragolangnfly || freebsd)

package net

/*
#include <netdb.h>
*/
import "C"

const cgolangAddrInfoFlags = (C.AI_CANONNAME | C.AI_V4MAPPED | C.AI_ALL) & C.AI_MASK
