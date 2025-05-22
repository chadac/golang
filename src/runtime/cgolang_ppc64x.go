// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build ppc64 || ppc64le

package runtime

// crosscall_ppc64 calls into the runtime to set up the registers the
// Golang runtime expects and so the symbol it calls needs to be exported
// for external linking to work.
//
//golang:cgolang_export_static _cgolang_reginit
