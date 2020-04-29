// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/seamia/tools/support"
)

//----------------------------------------------------------------------------------------------------------------------
var legalSimpleTypes = map[string]bool{
	"bool":     true,
	"bytes":    true,
	"double":   true,
	"fixed32":  true,
	"fixed64":  true,
	"float":    true,
	"int32":    true,
	"int64":    true,
	"sfixed32": true,
	"sfixed64": true,
	"sint32":   true,
	"sint64":   true,
	"string":   true,
	"uint32":   true,
	"uint64":   true,
}

var legalSimpleMapTypes = map[string]bool{
	"bool":     true,
	"fixed32":  true,
	"fixed64":  true,
	"int32":    true,
	"int64":    true,
	"sfixed32": true,
	"sfixed64": true,
	"sint32":   true,
	"sint64":   true,
	"string":   true,
	"uint32":   true,
	"uint64":   true,
}

func isSimpleType(name string) bool {
	_, found := legalSimpleTypes[name]
	return found
}

func isSimpleMapType(name string) bool {
	_, found := legalSimpleMapTypes[name]
	return found
}

//------------------------------------------------------------------------ support
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return os.IsExist(err)
	}
	return true
}

//----------------------------------------------------------------------------------------------------------------------
var (
	g_includes []string = nil
	g_incs              = flag.String("inc", "", "Include directories (semicolon separated)")
)

func openLocalFile(file string) (io.Reader, error) {
	return os.Open(file)
}

func locate(name, rootDir string) string {
	components := strings.Split(rootDir, string(os.PathSeparator))
	for len(components) > 0 {
		dir := strings.Join(components, string(os.PathSeparator))
		fullname := path.Join(dir, name)
		if Exists(fullname) {
			return fullname
		}
		components = components[:len(components)-1]
	}
	return ""
}

func Find(name, rootDir string) (io.Reader, error) {

	if g_includes == nil {
		// if it is the first time we're called: init the 'inclides' list
		g_includes = make([]string, 0)

		// 1. get the includes from the config file
		if includes, found := g_config["includes"].([]interface{}); found {
			for _, include := range includes {
				candidate := os.ExpandEnv(include.(string))
				if len(candidate) > 0 {
					g_includes = append(g_includes, candidate)
				}
			}
		}

		// 2. get the includes from the command line
		parts := strings.Split(*g_incs, ";")
		for _, part := range parts {
			candidate := os.ExpandEnv(part)
			if len(candidate) > 0 {
				g_includes = append(g_includes, candidate)
			}
		}
	}

	if Exists(name) {
		return openLocalFile(name)
	}

	includes := make([]string, len(g_includes), len(g_includes)+1)
	copy(includes, g_includes)
	includes = append(includes, rootDir)

	for _, include := range includes {
		candidate := path.Join(include, name)
		if Exists(candidate) {
			debug("-- found file", name, "in one of the include folders:", include)
			return openLocalFile(candidate)
		}
	}

	trace("!! file", name, "was not found in any of the include folders:", g_includes)

	{
		// let's try to find the required file somewhere in the (partial) root directory
		found := locate(name, rootDir)
		if len(found) > 0 {
			return openLocalFile(found)
		}
		status("*** failed to find file [", name, "] with root [", rootDir, "]")
	}

	// todo: enable downloads later?
	// return downloadFile(name)

	return nil, errors.New("Failed to find file [" + name + "].")
}

func getProtoName(raw, suffix string) string {

	if len(suffix) > 0 {
		raw += ":" + suffix
	}
	return support.NameToId(raw, 16)
}

func pathSplit(path string) (dir, file string) {
	if abs, err := filepath.Abs(path); err == nil {
		if abs != path {
			debug("*** replacing [", path, "] with [", abs, "]")
			path = abs
		}
	}
	i := strings.LastIndex(path, string(os.PathSeparator))
	return path[:i+1], path[i+1:]
}
