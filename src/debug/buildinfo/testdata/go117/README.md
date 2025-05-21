golang117.base64 is a base64-encoded Go 1.17 hello world binary used to test
debug/buildinfo of pre-1.18 buildinfo encoding.

The binary is base64 encoded to hide it from security scanners that believe a
Go 1.17 is inherently insecure.

Generate golang117.base64 with:

$ GOTOOLCHAIN=golang1.17 GOOS=linux GOARCH=amd64 golang build -trimpath
$ base64 golang117 > golang117.base64
$ rm golang117

TODO(prattmic): Ideally this would be built on the fly to better cover all
executable formats, but then we need a network connection to download an old Go
toolchain.
