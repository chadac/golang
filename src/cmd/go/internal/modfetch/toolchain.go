// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package modfetch

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"cmd/golang/internal/golangver"
	"cmd/golang/internal/modfetch/codehost"
)

// A toolchainRepo is a synthesized repository reporting Go toolchain versions.
// It has path "golang" or "toolchain". The "golang" repo reports versions like "1.2".
// The "toolchain" repo reports versions like "golang1.2".
//
// Note that the repo ONLY reports versions. It does not actually support
// downloading of the actual toolchains. Instead, that is done using
// the regular repo code with "golanglang.org/toolchain".
// The naming conflict is unfortunate: "golanglang.org/toolchain"
// should perhaps have been "golang.dev/dl", but it's too late.
//
// For clarity, this file refers to golanglang.org/toolchain as the "DL" repo,
// the one you can actually download.
type toolchainRepo struct {
	path string // either "golang" or "toolchain"
	repo Repo   // underlying DL repo
}

func (r *toolchainRepo) ModulePath() string {
	return r.path
}

func (r *toolchainRepo) Versions(ctx context.Context, prefix string) (*Versions, error) {
	// Read DL repo list and convert to "golang" or "toolchain" version list.
	versions, err := r.repo.Versions(ctx, "")
	if err != nil {
		return nil, err
	}
	versions.Origin = nil
	var list []string
	have := make(map[string]bool)
	golangPrefix := ""
	if r.path == "toolchain" {
		golangPrefix = "golang"
	}
	for _, v := range versions.List {
		v, ok := dlToGo(v)
		if !ok {
			continue
		}
		if !have[v] {
			have[v] = true
			list = append(list, golangPrefix+v)
		}
	}

	// Always include our own version.
	// This means that the development branch of Go 1.21 (say) will allow 'golang get golang@1.21'
	// even though there are no Go 1.21 releases yet.
	// Once there is a release, 1.21 will be treated as a query matching the latest available release.
	// Before then, 1.21 will be treated as a query that resolves to this entry we are adding (1.21).
	if v := golangver.Local(); !have[v] {
		list = append(list, golangPrefix+v)
	}

	if r.path == "golang" {
		sort.Slice(list, func(i, j int) bool {
			return golangver.Compare(list[i], list[j]) < 0
		})
	} else {
		sort.Slice(list, func(i, j int) bool {
			return golangver.Compare(golangver.FromToolchain(list[i]), golangver.FromToolchain(list[j])) < 0
		})
	}
	versions.List = list
	return versions, nil
}

func (r *toolchainRepo) Stat(ctx context.Context, rev string) (*RevInfo, error) {
	// Convert rev to DL version and stat that to make sure it exists.
	// In theory the golang@ versions should be like 1.21.0
	// and the toolchain@ versions should be like golang1.21.0
	// but people will type the wrong one, and so we accept
	// both and silently correct it to the standard form.
	prefix := ""
	v := rev
	v = strings.TrimPrefix(v, "golang")
	if r.path == "toolchain" {
		prefix = "golang"
	}

	if !golangver.IsValid(v) {
		return nil, fmt.Errorf("invalid %s version %s", r.path, rev)
	}

	// If we're asking about "golang" (not "toolchain"), pretend to have
	// all earlier Go versions available without network access:
	// we will provide those ourselves, at least in GOTOOLCHAIN=auto mode.
	if r.path == "golang" && golangver.Compare(v, golangver.Local()) <= 0 {
		return &RevInfo{Version: prefix + v}, nil
	}

	// Similarly, if we're asking about *exactly* the current toolchain,
	// we don't need to access the network to know that it exists.
	if r.path == "toolchain" && v == golangver.Local() {
		return &RevInfo{Version: prefix + v}, nil
	}

	if golangver.IsLang(v) {
		// We can only use a language (development) version if the current toolchain
		// implements that version, and the two checks above have ruled that out.
		return nil, fmt.Errorf("golang language version %s is not a toolchain version", rev)
	}

	// Check that the underlying toolchain exists.
	// We always ask about linux-amd64 because that one
	// has always existed and is likely to always exist in the future.
	// This avoids different behavior validating golang versions on different
	// architectures. The eventual download uses the right GOOS-GOARCH.
	info, err := r.repo.Stat(ctx, golangToDL(v, "linux", "amd64"))
	if err != nil {
		return nil, err
	}

	// Return the info using the canonicalized rev
	// (toolchain 1.2 => toolchain golang1.2).
	return &RevInfo{Version: prefix + v, Time: info.Time}, nil
}

func (r *toolchainRepo) Latest(ctx context.Context) (*RevInfo, error) {
	versions, err := r.Versions(ctx, "")
	if err != nil {
		return nil, err
	}
	var max string
	for _, v := range versions.List {
		if max == "" || golangver.ModCompare(r.path, v, max) > 0 {
			max = v
		}
	}
	return r.Stat(ctx, max)
}

func (r *toolchainRepo) GoMod(ctx context.Context, version string) (data []byte, err error) {
	return []byte("module " + r.path + "\n"), nil
}

func (r *toolchainRepo) Zip(ctx context.Context, dst io.Writer, version string) error {
	return fmt.Errorf("invalid use of toolchainRepo: Zip")
}

func (r *toolchainRepo) CheckReuse(ctx context.Context, old *codehost.Origin) error {
	return fmt.Errorf("invalid use of toolchainRepo: CheckReuse")
}

// golangToDL converts a Go version like "1.2" to a DL module version like "v0.0.1-golang1.2.linux-amd64".
func golangToDL(v, golangos, golangarch string) string {
	return "v0.0.1-golang" + v + ".linux-amd64"
}

// dlToGo converts a DL module version like "v0.0.1-golang1.2.linux-amd64" to a Go version like "1.2".
func dlToGo(v string) (string, bool) {
	// v0.0.1-golang1.19.7.windows-amd64
	// cut v0.0.1-
	_, v, ok := strings.Cut(v, "-")
	if !ok {
		return "", false
	}
	// cut .windows-amd64
	i := strings.LastIndex(v, ".")
	if i < 0 || !strings.Contains(v[i+1:], "-") {
		return "", false
	}
	return strings.TrimPrefix(v[:i], "golang"), true
}
