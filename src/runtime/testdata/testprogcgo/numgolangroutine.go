// Copyright 2017 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !plan9 && !windows
// +build !plan9,!windows

package main

/*
#include <stddef.h>
#include <pthread.h>

extern void CallbackNumGolangroutine();

static void* thread2(void* arg __attribute__ ((unused))) {
	CallbackNumGolangroutine();
	return NULL;
}

static void CheckNumGolangroutine() {
	pthread_t tid;
	pthread_create(&tid, NULL, thread2, NULL);
	pthread_join(tid, NULL);
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"strings"
)

var baseGolangroutines int

func init() {
	register("NumGolangroutine", NumGolangroutine)
}

func NumGolangroutine() {
	// Test that there are just the expected number of golangroutines
	// running. Specifically, test that the spare M's golangroutine
	// doesn't show up.
	if _, ok := checkNumGolangroutine("first", 1+baseGolangroutines); !ok {
		return
	}

	// Test that the golangroutine for a callback from C appears.
	if C.CheckNumGolangroutine(); !callbackok {
		return
	}

	// Make sure we're back to the initial golangroutines.
	if _, ok := checkNumGolangroutine("third", 1+baseGolangroutines); !ok {
		return
	}

	fmt.Println("OK")
}

func checkNumGolangroutine(label string, want int) (string, bool) {
	n := runtime.NumGolangroutine()
	if n != want {
		fmt.Printf("%s NumGolangroutine: want %d; golangt %d\n", label, want, n)
		return "", false
	}

	sbuf := make([]byte, 32<<10)
	sbuf = sbuf[:runtime.Stack(sbuf, true)]
	n = strings.Count(string(sbuf), "golangroutine ")
	if n != want {
		fmt.Printf("%s Stack: want %d; golangt %d:\n%s\n", label, want, n, sbuf)
		return "", false
	}
	return string(sbuf), true
}

var callbackok bool

//export CallbackNumGolangroutine
func CallbackNumGolangroutine() {
	stk, ok := checkNumGolangroutine("second", 2+baseGolangroutines)
	if !ok {
		return
	}
	if !strings.Contains(stk, "CallbackNumGolangroutine") {
		fmt.Printf("missing CallbackNumGolangroutine from stack:\n%s\n", stk)
		return
	}

	callbackok = true
}
