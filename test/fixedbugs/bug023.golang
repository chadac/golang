// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

type Type interface {
	TypeName() string;
}

type TInt struct {
}

// TInt
func (i *TInt) TypeName() string {
	return "int";
}


func main() {
	var t Type;
	t = nil;
	_ = t;
}

/*
bug023.golang:20: fatal error: naddr: const <Type>I{<TypeName>110(<_t117>{},<_o119>{},{});}
*/
