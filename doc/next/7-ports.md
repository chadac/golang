## Ports {#ports}

### Darwin

<!-- golang.dev/issue/69839 -->
As [announced](/doc/golang1.24#darwin) in the Go 1.24 release notes, Go 1.25 requires macOS 12 Monterey or later; support for previous versions has been discontinued.

### Windows

<!-- golang.dev/issue/71671 -->
Go 1.25 is the last release that contains the [broken](/doc/golang1.24#windows) 32-bit windows/arm port (`GOOS=windows` `GOARCH=arm`). It will be removed in Go 1.26.
