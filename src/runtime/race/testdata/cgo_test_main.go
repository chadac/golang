// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
int sync;

void Notify(void)
{
	__sync_fetch_and_add(&sync, 1);
}

void Wait(void)
{
	while(__sync_fetch_and_add(&sync, 0) == 0) {}
}
*/
import "C"

func main() {
	data := 0
	golang func() {
		data = 1
		C.Notify()
	}()
	C.Wait()
	_ = data
}
