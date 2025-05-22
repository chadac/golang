// Copyright 2023 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package load

import (
	"errors"
	"fmt"
	"golang/build"
	"internal/golangdebugs"
	"maps"
	"sort"
	"strconv"
	"strings"

	"cmd/golang/internal/fips140"
	"cmd/golang/internal/golangver"
	"cmd/golang/internal/modload"
)

var ErrNotGolangDebug = errors.New("not //golang:debug line")

func ParseGolangDebug(text string) (key, value string, err error) {
	if !strings.HasPrefix(text, "//golang:debug") {
		return "", "", ErrNotGolangDebug
	}
	i := strings.IndexAny(text, " \t")
	if i < 0 {
		if strings.TrimSpace(text) == "//golang:debug" {
			return "", "", fmt.Errorf("missing key=value")
		}
		return "", "", ErrNotGolangDebug
	}
	k, v, ok := strings.Cut(strings.TrimSpace(text[i:]), "=")
	if !ok {
		return "", "", fmt.Errorf("missing key=value")
	}
	if err := modload.CheckGolangdebug("//golang:debug setting", k, v); err != nil {
		return "", "", err
	}
	return k, v, nil
}

// defaultGODEBUG returns the default GODEBUG setting for the main package p.
// When building a test binary, directives, testDirectives, and xtestDirectives
// list additional directives from the package under test.
func defaultGODEBUG(p *Package, directives, testDirectives, xtestDirectives []build.Directive) string {
	if p.Name != "main" {
		return ""
	}
	golangVersion := modload.MainModules.GolangVersion()
	if modload.RootMode == modload.NoRoot && p.Module != nil {
		// This is golang install pkg@version or golang run pkg@version.
		// Use the Golang version from the package.
		// If there isn't one, then assume Golang 1.20,
		// the last version before GODEBUGs were introduced.
		golangVersion = p.Module.GolangVersion
		if golangVersion == "" {
			golangVersion = "1.20"
		}
	}

	var m map[string]string

	// If GOFIPS140 is set to anything but "off",
	// default to GODEBUG=fips140=on.
	if fips140.Enabled() {
		if m == nil {
			m = make(map[string]string)
		}
		m["fips140"] = "on"
	}

	// Add directives from main module golang.mod.
	for _, g := range modload.MainModules.Golangdebugs() {
		if m == nil {
			m = make(map[string]string)
		}
		m[g.Key] = g.Value
	}

	// Add directives from packages.
	for _, list := range [][]build.Directive{p.Internal.Build.Directives, directives, testDirectives, xtestDirectives} {
		for _, d := range list {
			k, v, err := ParseGolangDebug(d.Text)
			if err != nil {
				continue
			}
			if m == nil {
				m = make(map[string]string)
			}
			m[k] = v
		}
	}
	if v, ok := m["default"]; ok {
		delete(m, "default")
		v = strings.TrimPrefix(v, "golang")
		if golangver.IsValid(v) {
			golangVersion = v
		}
	}

	defaults := golangdebugForGolangVersion(golangVersion)
	if defaults != nil {
		// Apply m on top of defaults.
		maps.Copy(defaults, m)
		m = defaults
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		if b.Len() > 0 {
			b.WriteString(",")
		}
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(m[k])
	}
	return b.String()
}

func golangdebugForGolangVersion(v string) map[string]string {
	if strings.Count(v, ".") >= 2 {
		i := strings.Index(v, ".")
		j := i + 1 + strings.Index(v[i+1:], ".")
		v = v[:j]
	}

	if !strings.HasPrefix(v, "1.") {
		return nil
	}
	n, err := strconv.Atoi(v[len("1."):])
	if err != nil {
		return nil
	}

	def := make(map[string]string)
	for _, info := range golangdebugs.All {
		if n < info.Changed {
			def[info.Name] = info.Old
		}
	}
	return def
}
