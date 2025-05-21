// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

"use strict";

if (process.argv.length < 3) {
	console.error("usage: golang_js_wasm_exec [wasm binary] [arguments]");
	process.exit(1);
}

globalThis.require = require;
globalThis.fs = require("fs");
globalThis.path = require("path");
globalThis.TextEncoder = require("util").TextEncoder;
globalThis.TextDecoder = require("util").TextDecoder;

globalThis.performance ??= require("performance");

globalThis.crypto ??= require("crypto");

require("./wasm_exec");

const golang = new Go();
golang.argv = process.argv.slice(2);
golang.env = Object.assign({ TMPDIR: require("os").tmpdir() }, process.env);
golang.exit = process.exit;
WebAssembly.instantiate(fs.readFileSync(process.argv[2]), golang.importObject).then((result) => {
	process.on("exit", (code) => { // Node.js exits if no event handler is pending
		if (code === 0 && !golang.exited) {
			// deadlock, make Go print error and stack traces
			golang._pendingEvent = { id: 0 };
			golang._resume();
		}
	});
	return golang.run(result.instance);
}).catch((err) => {
	console.error(err);
	process.exit(1);
});
