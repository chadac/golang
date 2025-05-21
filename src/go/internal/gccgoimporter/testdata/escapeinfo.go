// Test case for escape info in export data. To compile and extract .golangx file:
// gccgolang -fgolang-optimize-allocs -c escapeinfo.golang
// objcopy -j .golang_export escapeinfo.o escapeinfo.golangx

package escapeinfo

type T struct{ data []byte }

func NewT(data []byte) *T {
	return &T{data}
}

func (*T) Read(p []byte) {}
