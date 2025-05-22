// errorcheck -lang=golang1.22

//golang:build golang1.21

// We need a line directive before the package clause,
// but don't change file name or position so that the
// error message appears at the right place.

//line issue67141.golang:10
package p

func _() {
	for range 10 { // ERROR "cannot range over 10"
	}
}
