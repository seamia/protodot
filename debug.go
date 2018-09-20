// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
)

const debugTrace = 2
const debugDebug = 4
const debugAlert = 8
const debugStatus = 16

const debugALL = debugStatus + debugAlert + debugDebug + debugTrace
const debugNone = 0
const debugNormal = debugStatus

var g_logWriter io.Writer
var g_debugLevel int = debugNormal

func alert(a ...interface{}) {
	if (g_debugLevel & debugAlert) == debugAlert {
		fmt.Println("ALERT!!!!", a)
		if g_logWriter != nil {
			fmt.Fprintln(g_logWriter, a...)
		}
	}
}

func assert(msg ...interface{}) {
	fmt.Print("ASSERT: ")
	fmt.Println(msg...)
	panic(1)
}

func trace(a ...interface{}) {
	if (g_debugLevel & debugTrace) == debugTrace {
		fmt.Println(a...)
		if g_logWriter != nil {
			// fmt.Fprintln(g_logWriter, a...)
		}
	}
}

func debug(a ...interface{}) {
	if (g_debugLevel & debugDebug) == debugDebug {
		fmt.Println(a...)

		if g_logWriter != nil {
			// fmt.Fprintln(g_logWriter, a...)
		}
	}
}

func status(a ...interface{}) {
	if (g_debugLevel & debugStatus) == debugStatus {
		fmt.Println(a...)
		if g_logWriter != nil {
			fmt.Fprintln(g_logWriter, a...)
		}
	}
}

func ignoring(a ...interface{}) {
	// debug(a...)
}

func unhandled(a ...interface{}) {
	// debug(a...)
}
