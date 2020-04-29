// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/seamia/tools/support"
)

func oneword(t string) string {
	return "oneword(" + t + ")"
}

func title(t string) string {
	return "title(" + t + ")"
}

func color(name string) string {
	if g_config != nil {
		if colors, found := g_config["colors"]; found {
			colorMap := colors.(map[string]interface{})
			if colorMap != nil {
				if maps2, found := colorMap[name]; found {
					name = maps2.(string)
				}
			}
		}
	}
	c, found := support.GetColor(name)
	if found != nil {
		fmt.Println("failed to resolve color name", name)
	}
	return c
}

func settings(key string) string {
	if g_config != nil {
		if colors, found := g_config["settings"]; found {
			settingsMap := colors.(map[string]interface{})
			if settingsMap != nil {
				if value, found := settingsMap[key]; found {
					return value.(string)
				}
			}
		}
	}

	fmt.Println("failed to resolve setting name", key)
	return "setting[" + key + "]"
}

// type FuncMap map[string]interface{}
var templFuncs = template.FuncMap{
	"lower":    strings.ToLower,
	"title":    title,
	"settings": settings,
	"color":    color,
	"oneword":  oneword,
}
