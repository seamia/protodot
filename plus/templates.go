// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plus

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/seamia/tools/assets"
	"github.com/seamia/tools/support"
)

var g_preloadedTemplates map[string]*template.Template

func loadExternals(name, tmplDir string) (string, error) {

	if len(tmplDir) > 0 {
		tmp := os.ExpandEnv(filepath.Join(tmplDir, name))
		if support.Exists(tmp) {
			data, err := ioutil.ReadFile(tmp)
			if err == nil {
				return string(data[:]), nil
			}
			return "", err
		}
	}

	name = os.ExpandEnv(name)
	if support.Exists(name) {
		data, err := ioutil.ReadFile(name)
		if err == nil {
			return string(data[:]), nil
		}
		return "", err
	}

	// last-ditch effort: let's look into the assets:
	if reader, err := assets.Open(name); err == nil {
		if raw, err := ioutil.ReadAll(reader); err == nil {
			return string(raw[:]), nil
		}
	}

	return "", errors.New("file [" + name + "] not found")
}

func oneLineText(raw string) string {
	bits := strings.Split(raw, "\n")
	result := ""
	for _, bit := range bits {
		result += strings.TrimSpace(bit)
	}
	return result
}

func resolveExternals(text, tmplDir string) (string, error) {
	if strings.HasPrefix(text, "file:") {
		bits := strings.Split(text, ":")
		fileName := bits[len(bits)-1]
		if text, err := loadExternals(fileName, tmplDir); err == nil {
			if len(bits) > 2 {
				options := bits[1]
				if options == "oneline" {
					text = oneLineText(text)
				}
			}
			return text, nil
		} else {
			return "", err
		}
	}
	return text, nil
}

func PreloadTemplates(config map[string]interface{}, funcs template.FuncMap, tmplDir string) error {

	g_preloadedTemplates = make(map[string]*template.Template)
	for name, data := range config {
		if text, err := resolveExternals(data.(string), tmplDir); err == nil {
			if tmpl, err := template.New(name).Funcs(funcs).Parse(text); err != nil {
				return err
			} else {
				g_preloadedTemplates[name] = tmpl
			}
		} else {
			return err
		}
	}
	return nil
}

func ApplyTemplate(name string, where io.Writer, payload interface{}) error {

	if g_preloadedTemplates != nil {
		if tmpl, present := g_preloadedTemplates[name]; present {
			return tmpl.Execute(where, payload)
		}
	}
	return errors.New("templates are not available")
}
