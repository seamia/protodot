// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"os"

	"github.com/seamia/tools/support"
)

type CreateOnWrite struct {
	name   string
	writer io.Writer
	err    error
}

func NewCreateOnWrite(name string) *CreateOnWrite {
	return &CreateOnWrite{name: name}
}

func (cow *CreateOnWrite) Write(p []byte) (n int, err error) {
	if cow.err != nil {
		return 0, cow.err
	}

	if cow.writer == nil {
		status("creating file:", cow.name)
		cow.writer, cow.err = os.Create(cow.name)
	}

	if cow.err != nil {
		return 0, cow.err
	}
	return cow.writer.Write(p)
}

// ----------------------------------------------------------------------------------------------------------------------
type ForkWriter struct {
	writers []io.Writer
}

func NewForkWriter(targets ...io.Writer) *ForkWriter {
	fork := ForkWriter{writers: targets}
	return &fork
}

func (fw *ForkWriter) AddWriter(target io.Writer) {
	if target != nil {
		if fw.writers == nil {
			fw.writers = make([]io.Writer, 0, 2)
		}
		fw.writers = append(fw.writers, target)
	}
}

func (fw *ForkWriter) Write(p []byte) (n int, err error) {
	for _, writer := range fw.writers {
		n, err = writer.Write(p) // ignoring (aka overwriting) 'intermediate' return values here
	}
	return
}

func createDirIfMissing(name string) {
	expanded := os.ExpandEnv(name)
	if len(expanded) > 0 {
		if !Exists(expanded) {
			status("creating missing direcory:", expanded)
			os.MkdirAll(expanded, 0755) // warning: dropping error on the floor here
		}
	}
}

func loadFileAsBytes(name string) []byte {
	if bytes, err := os.ReadFile(name); err == nil {
		return bytes
	} else {
		assert("failed to open/read file [", name, "], with error:", err)
	}
	return nil
}

func getTargetPath() string {
	var target string

	// 1. pull location from config file
	if dir, err := support.GetLocation(g_config, entryGenerated); err == nil {
		target = dir
	}

	// 2. use cli provided value (if present)
	if len(*g_generated) > 0 {
		target = *g_generated
	}

	return target
}
