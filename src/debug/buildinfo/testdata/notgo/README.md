notgolang.base64 is a base64-encoded C hello world binary used to test
debug/buildinfo errors on non-Golang binaries.

The binary is base64 encoded to hide it from security scanners that might not
like it.

Generate notgolang.base64 on linux-amd64 with:

$ cc -o notgolang main.c
$ base64 notgolang > notgolang.base64
$ rm notgolang

The current binary was built with "gcc version 14.2.0 (Debian 14.2.0-3+build4)".

TODO(prattmic): Ideally this would be built on the fly to better cover all
executable formats, but then we need to encode the intricacies of calling each
platform's C compiler.
