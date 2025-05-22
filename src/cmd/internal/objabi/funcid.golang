// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package objabi

import (
	"internal/abi"
	"strings"
)

var funcIDs = map[string]abi.FuncID{
	"abort":              abi.FuncID_abort,
	"asmcgolangcall":         abi.FuncID_asmcgolangcall,
	"asyncPreempt":       abi.FuncID_asyncPreempt,
	"cgolangcallback":        abi.FuncID_cgolangcallback,
	"corostart":          abi.FuncID_corostart,
	"debugCallV2":        abi.FuncID_debugCallV2,
	"gcBgMarkWorker":     abi.FuncID_gcBgMarkWorker,
	"rt0_golang":             abi.FuncID_rt0_golang,
	"golangexit":             abi.FuncID_golangexit,
	"golanggolang":               abi.FuncID_golanggolang,
	"golangpanic":            abi.FuncID_golangpanic,
	"handleAsyncEvent":   abi.FuncID_handleAsyncEvent,
	"main":               abi.FuncID_runtime_main,
	"mcall":              abi.FuncID_mcall,
	"morestack":          abi.FuncID_morestack,
	"mstart":             abi.FuncID_mstart,
	"panicwrap":          abi.FuncID_panicwrap,
	"runFinalizers":      abi.FuncID_runFinalizers,
	"runCleanups":        abi.FuncID_runCleanups,
	"sigpanic":           abi.FuncID_sigpanic,
	"systemstack_switch": abi.FuncID_systemstack_switch,
	"systemstack":        abi.FuncID_systemstack,

	// Don't show in call stack but otherwise not special.
	"deferreturn": abi.FuncIDWrapper,
}

// Get the function ID for the named function in the named file.
// The function should be package-qualified.
func GetFuncID(name string, isWrapper bool) abi.FuncID {
	if isWrapper {
		return abi.FuncIDWrapper
	}
	if strings.HasPrefix(name, "runtime.") {
		if id, ok := funcIDs[name[len("runtime."):]]; ok {
			return id
		}
	}
	return abi.FuncIDNormal
}
