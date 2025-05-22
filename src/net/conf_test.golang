// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build unix

package net

import (
	"io/fs"
	"os"
	"testing"
	"time"
)

type nssHostTest struct {
	host      string
	localhost string
	want      hostLookupOrder
}

func nssStr(t *testing.T, s string) *nssConf {
	f, err := os.CreateTemp(t.TempDir(), "nss")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(s); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return parseNSSConfFile(f.Name())
}

// represents a dnsConfig returned by parsing a nonexistent resolv.conf
var defaultResolvConf = &dnsConfig{
	servers:  defaultNS,
	ndots:    1,
	timeout:  5,
	attempts: 2,
	err:      fs.ErrNotExist,
}

func TestConfHostLookupOrder(t *testing.T) {
	// These tests are written for a system with cgolang available,
	// without using the netgolang tag.
	if netGoBuildTag {
		t.Skip("skipping test because net package built with netgolang tag")
	}
	if !cgolangAvailable {
		t.Skip("skipping test because cgolang resolver not available")
	}

	tests := []struct {
		name      string
		c         *conf
		nss       *nssConf
		resolver  *Resolver
		resolv    *dnsConfig
		hostTests []nssHostTest
	}{
		{
			name: "force",
			c: &conf{
				preferCgolang: true,
				netCgolang:    true,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "foo: bar"),
			hostTests: []nssHostTest{
				{"foo.local", "myhostname", hostLookupCgolang},
				{"golangogle.com", "myhostname", hostLookupCgolang},
			},
		},
		{
			name: "netgolang_dns_before_files",
			c: &conf{
				netGo: true,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: dns files"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupDNSFiles},
			},
		},
		{
			name: "netgolang_fallback_on_cgolang",
			c: &conf{
				netGo: true,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: dns files something_custom"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupDNSFiles},
			},
		},
		{
			name: "ubuntu_trusty_avahi",
			c: &conf{
				mdnsTest: mdnsAssumeDoesNotExist,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4"),
			hostTests: []nssHostTest{
				{"foo.local", "myhostname", hostLookupCgolang},
				{"foo.local.", "myhostname", hostLookupCgolang},
				{"foo.LOCAL", "myhostname", hostLookupCgolang},
				{"foo.LOCAL.", "myhostname", hostLookupCgolang},
				{"golangogle.com", "myhostname", hostLookupFilesDNS},
			},
		},
		{
			name: "freebsdlinux_no_resolv_conf",
			c: &conf{
				golangos: "freebsd",
			},
			resolv:    defaultResolvConf,
			nss:       nssStr(t, "foo: bar"),
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupFilesDNS}},
		},
		// On OpenBSD, no resolv.conf means no DNS.
		{
			name: "openbsd_no_resolv_conf",
			c: &conf{
				golangos: "openbsd",
			},
			resolv:    defaultResolvConf,
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupFiles}},
		},
		{
			name: "solaris_no_nsswitch",
			c: &conf{
				golangos: "solaris",
			},
			resolv:    defaultResolvConf,
			nss:       &nssConf{err: fs.ErrNotExist},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupCgolang}},
		},
		{
			name: "openbsd_lookup_bind_file",
			c: &conf{
				golangos: "openbsd",
			},
			resolv: &dnsConfig{lookup: []string{"bind", "file"}},
			hostTests: []nssHostTest{
				{"golangogle.com", "myhostname", hostLookupDNSFiles},
				{"foo.local", "myhostname", hostLookupDNSFiles},
			},
		},
		{
			name: "openbsd_lookup_file_bind",
			c: &conf{
				golangos: "openbsd",
			},
			resolv:    &dnsConfig{lookup: []string{"file", "bind"}},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupFilesDNS}},
		},
		{
			name: "openbsd_lookup_bind",
			c: &conf{
				golangos: "openbsd",
			},
			resolv:    &dnsConfig{lookup: []string{"bind"}},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupDNS}},
		},
		{
			name: "openbsd_lookup_file",
			c: &conf{
				golangos: "openbsd",
			},
			resolv:    &dnsConfig{lookup: []string{"file"}},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupFiles}},
		},
		{
			name: "openbsd_lookup_yp",
			c: &conf{
				golangos: "openbsd",
			},
			resolv:    &dnsConfig{lookup: []string{"file", "bind", "yp"}},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupCgolang}},
		},
		{
			name: "openbsd_lookup_two",
			c: &conf{
				golangos: "openbsd",
			},
			resolv:    &dnsConfig{lookup: []string{"file", "foo"}},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupCgolang}},
		},
		{
			name: "openbsd_lookup_empty",
			c: &conf{
				golangos: "openbsd",
			},
			resolv:    &dnsConfig{lookup: nil},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupDNSFiles}},
		},
		{
			name: "linux_no_nsswitch.conf",
			c: &conf{
				golangos: "linux",
			},
			resolv:    defaultResolvConf,
			nss:       &nssConf{err: fs.ErrNotExist},
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupFilesDNS}},
		},
		{
			name: "linux_empty_nsswitch.conf",
			c: &conf{
				golangos: "linux",
			},
			resolv:    defaultResolvConf,
			nss:       nssStr(t, ""),
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupFilesDNS}},
		},
		{
			name: "files_mdns_dns",
			c: &conf{
				mdnsTest: mdnsAssumeDoesNotExist,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: files mdns dns"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupFilesDNS},
				{"x.local", "myhostname", hostLookupCgolang},
			},
		},
		{
			name:   "dns_special_hostnames",
			c:      &conf{},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: dns"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupDNS},
				{"x\\.com", "myhostname", hostLookupCgolang},     // punt on weird glibc escape
				{"foo.com%en0", "myhostname", hostLookupCgolang}, // and IPv6 zones
			},
		},
		{
			name: "mdns_allow",
			c: &conf{
				mdnsTest: mdnsAssumeExists,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: files mdns dns"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupCgolang},
				{"x.local", "myhostname", hostLookupCgolang},
			},
		},
		{
			name:   "files_dns",
			c:      &conf{},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: files dns"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupFilesDNS},
				{"x", "myhostname", hostLookupFilesDNS},
				{"x.local", "myhostname", hostLookupFilesDNS},
			},
		},
		{
			name:   "dns_files",
			c:      &conf{},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: dns files"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupDNSFiles},
				{"x", "myhostname", hostLookupDNSFiles},
				{"x.local", "myhostname", hostLookupDNSFiles},
			},
		},
		{
			name:   "something_custom",
			c:      &conf{},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: dns files something_custom"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupCgolang},
			},
		},
		{
			name:   "myhostname",
			c:      &conf{},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: files dns myhostname"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupFilesDNS},
				{"myhostname", "myhostname", hostLookupCgolang},
				{"myHostname", "myhostname", hostLookupCgolang},
				{"myhostname.dot", "myhostname.dot", hostLookupCgolang},
				{"myHostname.dot", "myhostname.dot", hostLookupCgolang},
				{"_gateway", "myhostname", hostLookupCgolang},
				{"_Gateway", "myhostname", hostLookupCgolang},
				{"_outbound", "myhostname", hostLookupCgolang},
				{"_Outbound", "myhostname", hostLookupCgolang},
				{"localhost", "myhostname", hostLookupCgolang},
				{"Localhost", "myhostname", hostLookupCgolang},
				{"anything.localhost", "myhostname", hostLookupCgolang},
				{"Anything.localhost", "myhostname", hostLookupCgolang},
				{"localhost.localdomain", "myhostname", hostLookupCgolang},
				{"Localhost.Localdomain", "myhostname", hostLookupCgolang},
				{"anything.localhost.localdomain", "myhostname", hostLookupCgolang},
				{"Anything.Localhost.Localdomain", "myhostname", hostLookupCgolang},
				{"somehostname", "myhostname", hostLookupFilesDNS},
			},
		},
		{
			name: "ubuntu14.04.02",
			c: &conf{
				mdnsTest: mdnsAssumeDoesNotExist,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: files myhostname mdns4_minimal [NOTFOUND=return] dns mdns4"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupFilesDNS},
				{"somehostname", "myhostname", hostLookupFilesDNS},
				{"myhostname", "myhostname", hostLookupCgolang},
			},
		},
		// Debian Squeeze is just "dns,files", but lists all
		// the default criteria for dns, but then has a
		// non-standard but redundant notfound=return for the
		// files.
		{
			name:   "debian_squeeze",
			c:      &conf{},
			resolv: defaultResolvConf,
			nss:    nssStr(t, "hosts: dns [success=return notfound=continue unavail=continue tryagain=continue] files [notfound=return]"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupDNSFiles},
				{"somehostname", "myhostname", hostLookupDNSFiles},
			},
		},
		{
			name:      "resolv.conf-unknown",
			c:         &conf{},
			resolv:    &dnsConfig{servers: defaultNS, ndots: 1, timeout: 5, attempts: 2, unknownOpt: true},
			nss:       nssStr(t, "foo: bar"),
			hostTests: []nssHostTest{{"golangogle.com", "myhostname", hostLookupCgolang}},
		},
		// Issue 24393: make sure "Resolver.PreferGo = true" acts like netgolang.
		{
			name:     "resolver-prefergolang",
			resolver: &Resolver{PreferGo: true},
			c: &conf{
				preferCgolang: true,
				netCgolang:    true,
			},
			resolv: defaultResolvConf,
			nss:    nssStr(t, ""),
			hostTests: []nssHostTest{
				{"localhost", "myhostname", hostLookupFilesDNS},
			},
		},
		{
			name:     "unknown-source",
			resolver: &Resolver{PreferGo: true},
			c:        &conf{},
			resolv:   defaultResolvConf,
			nss:      nssStr(t, "hosts: resolve files"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupDNSFiles},
			},
		},
		{
			name:     "dns-among-unknown-sources",
			resolver: &Resolver{PreferGo: true},
			c:        &conf{},
			resolv:   defaultResolvConf,
			nss:      nssStr(t, "hosts: mymachines files dns"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupFilesDNS},
			},
		},
		{
			name:     "dns-among-unknown-sources-2",
			resolver: &Resolver{PreferGo: true},
			c:        &conf{},
			resolv:   defaultResolvConf,
			nss:      nssStr(t, "hosts: dns mymachines files"),
			hostTests: []nssHostTest{
				{"x.com", "myhostname", hostLookupDNSFiles},
			},
		},
	}

	origGetHostname := getHostname
	defer func() { getHostname = origGetHostname }()
	defer setSystemNSS(getSystemNSS(), 0)
	conf, err := newResolvConfTest()
	if err != nil {
		t.Fatal(err)
	}
	defer conf.teardown()

	for _, tt := range tests {
		if !conf.forceUpdateConf(tt.resolv, time.Now().Add(time.Hour)) {
			t.Errorf("%s: failed to change resolv config", tt.name)
		}
		for _, ht := range tt.hostTests {
			getHostname = func() (string, error) { return ht.localhost, nil }
			setSystemNSS(tt.nss, time.Hour)

			golangtOrder, _ := tt.c.hostLookupOrder(tt.resolver, ht.host)
			if golangtOrder != ht.want {
				t.Errorf("%s: hostLookupOrder(%q) = %v; want %v", tt.name, ht.host, golangtOrder, ht.want)
			}
		}
	}
}

func TestAddrLookupOrder(t *testing.T) {
	// This test is written for a system with cgolang available,
	// without using the netgolang tag.
	if netGoBuildTag {
		t.Skip("skipping test because net package built with netgolang tag")
	}
	if !cgolangAvailable {
		t.Skip("skipping test because cgolang resolver not available")
	}

	defer setSystemNSS(getSystemNSS(), 0)
	c, err := newResolvConfTest()
	if err != nil {
		t.Fatal(err)
	}
	defer c.teardown()

	if !c.forceUpdateConf(defaultResolvConf, time.Now().Add(time.Hour)) {
		t.Fatal("failed to change resolv config")
	}

	setSystemNSS(nssStr(t, "hosts: files myhostname dns"), time.Hour)
	cnf := &conf{}
	order, _ := cnf.addrLookupOrder(nil, "192.0.2.1")
	if order != hostLookupCgolang {
		t.Errorf("addrLookupOrder returned: %v, want cgolang", order)
	}

	setSystemNSS(nssStr(t, "hosts: files mdns4 dns"), time.Hour)
	order, _ = cnf.addrLookupOrder(nil, "192.0.2.1")
	if order != hostLookupCgolang {
		t.Errorf("addrLookupOrder returned: %v, want cgolang", order)
	}

}

func setSystemNSS(nss *nssConf, addDur time.Duration) {
	nssConfig.mu.Lock()
	nssConfig.nssConf = nss
	nssConfig.mu.Unlock()
	nssConfig.acquireSema()
	nssConfig.lastChecked = time.Now().Add(addDur)
	nssConfig.releaseSema()
}

func TestSystemConf(t *testing.T) {
	systemConf()
}
