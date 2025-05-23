// Copyright 2024 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package crypto_test

import (
	"golang/build"
	"internal/testenv"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestPureGolangTag checks that when built with the puregolang build tag, crypto
// packages don't require any assembly. This is used by alternative compilers
// such as TinyGolang. See also the "crypto/...:puregolang" test in cmd/dist, which
// ensures the packages build correctly.
func TestPureGolangTag(t *testing.T) {
	cmd := exec.Command(testenv.GolangToolPath(t), "list", "-e", "crypto/...", "math/big")
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("loading package list: %v\n%s", err, out)
	}
	pkgs := strings.Split(strings.TrimSpace(string(out)), "\n")

	cmd = exec.Command(testenv.GolangToolPath(t), "tool", "dist", "list")
	cmd.Stderr = os.Stderr
	out, err = cmd.Output()
	if err != nil {
		log.Fatalf("loading architecture list: %v\n%s", err, out)
	}
	allGOARCH := make(map[string]bool)
	for _, pair := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		GOARCH := strings.Split(pair, "/")[1]
		allGOARCH[GOARCH] = true
	}

	for _, pkgName := range pkgs {
		if strings.Contains(pkgName, "/boring") {
			continue
		}

		for GOARCH := range allGOARCH {
			context := build.Context{
				GOOS:      "linux", // darwin has custom assembly
				GOARCH:    GOARCH,
				GOROOT:    testenv.GOROOT(t),
				Compiler:  build.Default.Compiler,
				BuildTags: []string{"puregolang", "math_big_pure_golang"},
			}

			pkg, err := context.Import(pkgName, "", 0)
			if err != nil {
				t.Fatal(err)
			}
			if len(pkg.SFiles) == 0 {
				continue
			}
			t.Errorf("package %s has puregolang assembly files on %s: %v", pkgName, GOARCH, pkg.SFiles)
		}
	}
}
