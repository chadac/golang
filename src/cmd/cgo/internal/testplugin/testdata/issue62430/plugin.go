package main

import (
	"unicode"
)

func F(s string) *unicode.RangeTable {
	return unicode.Categolangries[s]
}

func main() {}
