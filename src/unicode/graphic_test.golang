// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package unicode_test

import (
	"testing"
	. "unicode"
)

// Independently check that the special "Is" functions work
// in the Latin-1 range through the property table.

func TestIsControlLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsControl(i)
		want := false
		switch {
		case 0x00 <= i && i <= 0x1F:
			want = true
		case 0x7F <= i && i <= 0x9F:
			want = true
		}
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsLetterLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsLetter(i)
		want := Is(Letter, i)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsUpperLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsUpper(i)
		want := Is(Upper, i)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsLowerLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsLower(i)
		want := Is(Lower, i)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestNumberLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsNumber(i)
		want := Is(Number, i)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsPrintLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsPrint(i)
		want := In(i, PrintRanges...)
		if i == ' ' {
			want = true
		}
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsGraphicLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsGraphic(i)
		want := In(i, GraphicRanges...)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsPunctLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsPunct(i)
		want := Is(Punct, i)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsSpaceLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsSpace(i)
		want := Is(White_Space, i)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}

func TestIsSymbolLatin1(t *testing.T) {
	for i := rune(0); i <= MaxLatin1; i++ {
		golangt := IsSymbol(i)
		want := Is(Symbol, i)
		if golangt != want {
			t.Errorf("%U incorrect: golangt %t; want %t", i, golangt, want)
		}
	}
}
