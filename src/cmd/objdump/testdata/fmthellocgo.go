package main

import "fmt"
import "C"

func main() {
	Println("hello, world")
	if flag {
//line fmthello.golang:999999
		Println("bad line")
		for {
		}
	}
}

//golang:noinline
func Println(s string) {
	fmt.Println(s)
}

var flag bool
