// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"

	"github.com/seamia/tools/support"
)

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
				if output, e := exec.Command(graphviz, "-Tsvg", src).Output(); e == nil {
					if err := os.WriteFile(svgPath, output, 0755); err != nil {
						status("error on write", err)
						svgPath = ""
					}
				} else {
					status("error on exec", e)
					svgPath = ""
				}
			}

			if png {
				status("generating .png file")
				if output, e := exec.Command(graphviz, "-Tpng", src).Output(); e == nil {
					if err := os.WriteFile(pngPath, output, 0755); err != nil {
						status("error on write", err)
						pngPath = ""
					}
				} else {
					status("error on exec", e)
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
