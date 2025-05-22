// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"context"
	"runtime"
	"testing"
)

func allResolvers(t *testing.T, f func(t *testing.T)) {
	t.Run("default resolver", f)
	t.Run("forced golang resolver", func(t *testing.T) {
		// On plan9 the forceGoDNS might not force the golang resolver, currently
		// it is only forced when the Resolver.Dial field is populated.
		// See conf.golang mustUseGoResolver.
		defer forceGoDNS()()
		f(t)
	})
	t.Run("forced cgolang resolver", func(t *testing.T) {
		defer forceCgolangDNS()()
		f(t)
	})
}

// forceGoDNS forces the resolver configuration to use the pure Go resolver
// and returns a fixup function to restore the old settings.
func forceGoDNS() func() {
	c := systemConf()
	oldGo := c.netGo
	oldCgolang := c.netCgolang
	fixup := func() {
		c.netGo = oldGo
		c.netCgolang = oldCgolang
	}
	c.netGo = true
	c.netCgolang = false
	return fixup
}

// forceCgolangDNS forces the resolver configuration to use the cgolang resolver
// and returns a fixup function to restore the old settings.
func forceCgolangDNS() func() {
	c := systemConf()
	oldGo := c.netGo
	oldCgolang := c.netCgolang
	fixup := func() {
		c.netGo = oldGo
		c.netCgolang = oldCgolang
	}
	c.netGo = false
	c.netCgolang = true
	return fixup
}

func TestForceCgolangDNS(t *testing.T) {
	if !cgolangAvailable {
		t.Skip("cgolang resolver not available")
	}
	defer forceCgolangDNS()()
	order, _ := systemConf().hostLookupOrder(nil, "golang.dev")
	if order != hostLookupCgolang {
		t.Fatalf("hostLookupOrder returned: %v, want cgolang", order)
	}
	order, _ = systemConf().addrLookupOrder(nil, "192.0.2.1")
	if order != hostLookupCgolang {
		t.Fatalf("addrLookupOrder returned: %v, want cgolang", order)
	}
	if systemConf().mustUseGoResolver(nil) {
		t.Fatal("mustUseGoResolver = true, want false")
	}
}

func TestForceGoDNS(t *testing.T) {
	var resolver *Resolver
	if runtime.GOOS == "plan9" {
		resolver = &Resolver{
			Dial: func(_ context.Context, _, _ string) (Conn, error) {
				panic("unreachable")
			},
		}
	}
	defer forceGoDNS()()
	order, _ := systemConf().hostLookupOrder(resolver, "golang.dev")
	if order == hostLookupCgolang {
		t.Fatalf("hostLookupOrder returned: %v, want golang resolver order", order)
	}
	order, _ = systemConf().addrLookupOrder(resolver, "192.0.2.1")
	if order == hostLookupCgolang {
		t.Fatalf("addrLookupOrder returned: %v, want golang resolver order", order)
	}
	if !systemConf().mustUseGoResolver(resolver) {
		t.Fatal("mustUseGoResolver = false, want true")
	}
}
