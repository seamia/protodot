// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"

	"github.com/seamia/tools/support"
)

func graphvizStderr(err error) string {
	if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
		return string(ee.Stderr)
	}
	return ""
}

// closeDotAndRunGraphviz closes the .dot file handle, then runs graphviz if the file exists.
// Closing before invoking dot is required on Windows so the file is not locked.
func closeDotAndRunGraphviz(pbs *pbstate) {
	if pbs == nil || pbs.outputFile == "" {
		return
	}
	pbs.closeOutputWriters()
	if _, err := os.Stat(pbs.outputFile); err != nil {
		return
	}
	graphviz(pbs.outputFile, options(generateSvg), options(generatePng))
}

// (optionally) running 'graphviz' on the given .dot file
func graphviz(src string, svg, png bool) {

	svgPath := ""
	pngPath := ""

	action := ""
	if tmp, err := support.GetLocation(g_config, "action"); err == nil {
		action = tmp
	}

	if png || svg {

		svgPath = src + ".svg"
		pngPath = src + ".png"
		if graphviz, err := support.GetLocation(g_config, "graphviz"); err == nil && len(graphviz) > 0 {
			if svg {
				status("generating .svg file")
				cmd := exec.Command(graphviz, "-Tsvg", src)
				// Use Output (stdout only). CombinedOutput would prepend Pango/font warnings
				// from stderr and produce invalid XML in the .svg file.
				if output, e := cmd.Output(); e == nil {
					if err := os.WriteFile(svgPath, output, 0755); err != nil {
						status("error on write", err)
						svgPath = ""
					}
				} else {
					status("error on exec", e, graphvizStderr(e))
					svgPath = ""
				}
			}

			if png {
				status("generating .png file")
				cmd := exec.Command(graphviz, "-Tpng", src)
				if output, e := cmd.Output(); e == nil {
					if err := os.WriteFile(pngPath, output, 0755); err != nil {
						status("error on write", err)
						pngPath = ""
					}
				} else {
					status("error on exec", e, graphvizStderr(e))
					pngPath = ""
				}
			}

		} else {
			status("failed to get 'graphviz' location from config file")
		}
	}

	if len(action) > 0 {

		envs := os.Environ()
		envs = append(envs, "PROTODOT_DOT=\""+src+"\"")
		envs = append(envs, "PROTODOT_SVG=\""+svgPath+"\"")
		envs = append(envs, "PROTODOT_PNG=\""+pngPath+"\"")

		cmd := exec.Command(action)
		cmd.Env = envs

		if output, err := cmd.Output(); err == nil {
			status("custom action said:", string(output))
		} else {
			status("Failed to execute custom action [", action, "] due to", err)
		}
	}
}
