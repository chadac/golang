// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !android && cgolang && !netgolang

package net

/*
#include <netdb.h>
*/
import "C"

// NOTE(rsc): In theory there are approximately balanced
// arguments for and against including AI_ADDRCONFIG
// in the flags (it includes IPv4 results only on IPv4 systems,
// and similarly for IPv6), but in practice setting it causes
// getaddrinfo to return the wrong canonical name on Linux.
// So definitely leave it out.
const cgolangAddrInfoFlags = C.AI_CANONNAME | C.AI_V4MAPPED | C.AI_ALL
