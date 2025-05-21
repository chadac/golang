// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package syscall

//golang:wasmimport wasi_snapshot_preview1 proc_exit
func ProcExit(code int32)
