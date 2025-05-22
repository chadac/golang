// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Tests a many interesting cases (network, syscalls, a little GC, busy golangroutines,
// blocked golangroutines, LockOSThread, pipes, and GOMAXPROCS).

//golang:build ignore

package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(8))
	outerDone := make(chan bool)

	golang func() {
		defer func() {
			outerDone <- true
		}()

		var wg sync.WaitGroup
		done := make(chan bool)

		wg.Add(1)
		golang func() {
			<-done
			wg.Done()
		}()

		rp, wp, err := os.Pipe()
		if err != nil {
			log.Fatalf("failed to create pipe: %v", err)
			return
		}
		defer func() {
			rp.Close()
			wp.Close()
		}()
		wg.Add(1)
		golang func() {
			var tmp [1]byte
			rp.Read(tmp[:])
			<-done
			wg.Done()
		}()
		time.Sleep(time.Millisecond)

		golang func() {
			runtime.LockOSThread()
			for {
				select {
				case <-done:
					return
				default:
					runtime.Gosched()
				}
			}
		}()

		runtime.GC()
		// Trigger GC from malloc.
		n := 512
		for i := 0; i < n; i++ {
			_ = make([]byte, 1<<20)
		}

		// Create a bunch of busy golangroutines to load all Ps.
		for p := 0; p < 10; p++ {
			wg.Add(1)
			golang func() {
				// Do something useful.
				tmp := make([]byte, 1<<16)
				for i := range tmp {
					tmp[i]++
				}
				_ = tmp
				<-done
				wg.Done()
			}()
		}

		// Block in syscall.
		wg.Add(1)
		golang func() {
			var tmp [1]byte
			rp.Read(tmp[:])
			<-done
			wg.Done()
		}()

		runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))

		// Test timers.
		timerDone := make(chan bool)
		golang func() {
			time.Sleep(time.Millisecond)
			timerDone <- true
		}()
		<-timerDone

		// A bit of network.
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("listen failed: %v", err)
			return
		}
		defer ln.Close()
		golang func() {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			time.Sleep(time.Millisecond)
			var buf [1]byte
			c.Write(buf[:])
			c.Close()
		}()
		c, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			log.Fatalf("dial failed: %v", err)
			return
		}
		var tmp [1]byte
		c.Read(tmp[:])
		c.Close()

		golang func() {
			runtime.Gosched()
			select {}
		}()

		// Unblock helper golangroutines and wait them to finish.
		wp.Write(tmp[:])
		wp.Write(tmp[:])
		close(done)
		wg.Wait()
	}()

	const iters = 5
	for i := 0; i < iters; i++ {
		var w io.Writer
		if i == iters-1 {
			w = os.Stdout
		} else {
			w = new(bytes.Buffer)
		}
		if err := trace.Start(w); err != nil {
			log.Fatalf("failed to start tracing: %v", err)
		}
		time.Sleep(time.Millisecond)
		trace.Stop()
	}
	<-outerDone
}
