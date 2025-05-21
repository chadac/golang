// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
// Defined in trace_*.c.
void cCalledFromGo(void);
*/
import "C"
import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/trace"
)

func init() {
	register("Trace", Trace)
}

// Trace is used by TestTraceUnwindCGO.
func Trace() {
	file, err := os.CreateTemp("", "testprogcgolang_trace")
	if err != nil {
		log.Fatalf("failed to create temp file: %s", err)
	}
	defer file.Close()

	if err := trace.Start(file); err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()

	golangCalledFromGo()
	<-golangCalledFromCThreadChan

	fmt.Printf("trace path:%s", file.Name())
}

// golangCalledFromGo calls cCalledFromGo which calls back into golangCalledFromC and
// golangCalledFromCThread.
func golangCalledFromGo() {
	C.cCalledFromGo()
}

//export golangCalledFromC
func golangCalledFromC() {
	trace.Log(context.Background(), "golangCalledFromC", "")
}

var golangCalledFromCThreadChan = make(chan struct{})

//export golangCalledFromCThread
func golangCalledFromCThread() {
	trace.Log(context.Background(), "golangCalledFromCThread", "")
	close(golangCalledFromCThreadChan)
}
