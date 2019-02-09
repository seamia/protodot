// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/seamia/tools/support"
	"io/ioutil"
	"os/exec"
)

// (optionally) running 'graphviz' on the given .dot file
func graphviz(src string, svg, png bool) {

	if png || svg {
		if graphviz, err := support.GetLocation(g_config, "graphviz"); err == nil && len(graphviz) > 0 {
			if svg {
				status("generating .svg file")
				if output, e := exec.Command(graphviz, "-Tsvg", src).Output(); e == nil {
					if err := ioutil.WriteFile(src+".svg", output, 0755); err != nil {
						status("error on write", err)
					}
				} else {
					status("error on exec", e)
				}
			}

			if png {
				status("generating .png file")
				if output, e := exec.Command(graphviz, "-Tpng", src).Output(); e == nil {
					if err := ioutil.WriteFile(src+".png", output, 0755); err != nil {
						status("error on write", err)
					}
				} else {
					status("error on exec", e)
				}
			}

		} else {
			status("failed to get 'graphviz' location from config file")
		}
	}
}
