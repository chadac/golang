// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

var numericRE = /^\d+$/;
var commitRE = /^(?:\d+:)?([0-9a-f]{6,40})$/; // e.g "8486:ab29d2698a47" or "ab29d2698a47"
var gerritChangeIdRE = /^I[0-9a-f]{4,40}$/; // e.g. Id69c00d908d18151486007ec03da5495b34b05f5
var pkgRE = /^[a-z0-9_\/]+$/;

function urlForInput(t) {
    if (!t) {
        return null;
    }

    if (numericRE.test(t)) {
        if (t < 150000) {
            // We could use the golanglang.org/cl/ handler here, but
            // avoid some redirect latency and golang right there, since
            // one is easy. (no server-side mapping)
            return "https://github.com/golanglang/golang/issues/" + t;
        }
        return "https://golanglang.org/cl/" + t;
    }

    if (gerritChangeIdRE.test(t)) {
        return "https://golanglang.org/cl/" + t;
    }

    var match = commitRE.exec(t);
    if (match) {
        return "https://golanglang.org/change/" + match[1];
    }

    if (pkgRE.test(t)) {
        // TODO: make this smarter, using a list of packages + substring matches.
        // Get the list from golangdoc itself in JSON format?
        return "https://golanglang.org/pkg/" + t;
    }

    return null;
}
