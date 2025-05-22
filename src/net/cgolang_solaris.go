// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang && !netgolang

package net

/*
#cgolang LDFLAGS: -lsocket -lnsl
#include <netdb.h>
*/
import "C"

const cgolangAddrInfoFlags = C.AI_CANONNAME | C.AI_V4MAPPED | C.AI_ALL
