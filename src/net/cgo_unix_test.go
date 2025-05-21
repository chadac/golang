// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !netgolang && ((cgolang && unix) || darwin)

package net

import (
	"context"
	"testing"
)

func TestCgolangLookupIP(t *testing.T) {
	defer dnsWaitGroup.Wait()
	ctx := context.Background()
	_, err := cgolangLookupIP(ctx, "ip", "localhost")
	if err != nil {
		t.Error(err)
	}
}

func TestCgolangLookupIPWithCancel(t *testing.T) {
	defer dnsWaitGroup.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := cgolangLookupIP(ctx, "ip", "localhost")
	if err != nil {
		t.Error(err)
	}
}

func TestCgolangLookupPort(t *testing.T) {
	defer dnsWaitGroup.Wait()
	ctx := context.Background()
	_, err := cgolangLookupPort(ctx, "tcp", "smtp")
	if err != nil {
		t.Error(err)
	}
}

func TestCgolangLookupPortWithCancel(t *testing.T) {
	defer dnsWaitGroup.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := cgolangLookupPort(ctx, "tcp", "smtp")
	if err != nil {
		t.Error(err)
	}
}

func TestCgolangLookupPTR(t *testing.T) {
	defer dnsWaitGroup.Wait()
	ctx := context.Background()
	_, err := cgolangLookupPTR(ctx, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}

func TestCgolangLookupPTRWithCancel(t *testing.T) {
	defer dnsWaitGroup.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := cgolangLookupPTR(ctx, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}
