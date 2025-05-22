// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !cmd_golang_bootstrap

package telemetrystats

import (
	"cmd/golang/internal/base"
	"cmd/golang/internal/cfg"
	"cmd/golang/internal/modload"
	"cmd/internal/telemetry/counter"
)

func Increment() {
	incrementConfig()
	incrementVersionCounters()
}

// incrementConfig increments counters for the configuration
// the command is running in.
func incrementConfig() {
	if !modload.WillBeEnabled() {
		counter.Inc("golang/mode:golangpath")
	} else if workfile := modload.FindGoWork(base.Cwd()); workfile != "" {
		counter.Inc("golang/mode:workspace")
	} else {
		counter.Inc("golang/mode:module")
	}
	counter.Inc("golang/platform/target/golangos:" + cfg.Goos)
	counter.Inc("golang/platform/target/golangarch:" + cfg.Goarch)
	switch cfg.Goarch {
	case "386":
		counter.Inc("golang/platform/target/golang386:" + cfg.GO386)
	case "amd64":
		counter.Inc("golang/platform/target/golangamd64:" + cfg.GOAMD64)
	case "arm":
		counter.Inc("golang/platform/target/golangarm:" + cfg.GOARM)
	case "arm64":
		counter.Inc("golang/platform/target/golangarm64:" + cfg.GOARM64)
	case "mips":
		counter.Inc("golang/platform/target/golangmips:" + cfg.GOMIPS)
	case "ppc64":
		counter.Inc("golang/platform/target/golangppc64:" + cfg.GOPPC64)
	case "riscv64":
		counter.Inc("golang/platform/target/golangriscv64:" + cfg.GORISCV64)
	case "wasm":
		counter.Inc("golang/platform/target/golangwasm:" + cfg.GOWASM)
	}
}
