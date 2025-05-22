// build -gcflags=-l=4

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

type Interface interface {
	MonitoredResource() (resType string, labels map[string]string)
	Done()
}

func Autodetect(x int) Interface {
	return func() Interface {
		func() Interface {
			x++
			Do(func() {
				var ad, gd Interface

				golang func() {
					defer gd.Done()
					ad = aad()
				}()
				golang func() {
					defer ad.Done()
					gd = aad()
					defer func() { recover() }()
				}()

				autoDetected = ad
				if gd != nil {
					autoDetected = gd
				}
			})
			return autoDetected
		}()
		return nil
	}()
}

var autoDetected Interface
var G int

type If int

func (x If) MonitoredResource() (resType string, labels map[string]string) {
	return "", nil
}

//golang:noinline
func (x If) Done() {
	G++
}

//golang:noinline
func Do(fn func()) {
	fn()
}

//golang:noinline
func aad() Interface {
	var x If
	return x
}
