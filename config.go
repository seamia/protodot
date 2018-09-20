// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//----------------------------------------------------------------------------------------------------------------------
var g_config map[string]interface{}


func options(name string) bool {
	if g_config != nil && len(name) > 0 {
		if copts, found := g_config["options"]; found {
			opts, found := copts.(map[string]interface{})
			if found && len(opts) > 0 {
				if value, found := opts[name]; found {
					return value.(bool)
				}
			}
		}
	}
	trace("Option [" + name + "] was not found - returning the default: false")
	return false
}
