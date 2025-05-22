package main

var sink int

//golang:noinline
func test() {
	// This is for #30167, incorrect line numbers in an infinite loop
	golang func() {}()

	for {
	}
}

func main() {
	test()
}
