// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"errors"
	"github.com/seamia/tools/support"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func downloadFromUrl(url, filename string) (io.Reader, error) {

	if downloads,err := support.GetLocation(g_config, "downloads"); err == nil && len(downloads) > 0 {
		if len(filename) == 0 {
			tokens := strings.Split(url, "/")
			filename = path.Join(downloads, tokens[len(tokens)-1])
			trace("Downloading", url, "to", filename)
		} else {
			filename = path.Join(downloads, filename)
		}
	} else {
		alert("the 'download' location is not set")
		if err == nil {
			err = errors.New("the 'download' location is empty")
		}
		return nil, err
	}

	if Exists(filename) {
		trace("file already exists. ", url, " maps to ", filename)
		return os.Open(filename)	// return from cache
	}

	response, err := http.Get(url)
	if err != nil {
		alert("Error while downloading", url, "-", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		alert("got invalid status.code", response.StatusCode, "while downloading", url)
		return nil, errors.New("failed to download: "+response.Status)
	}

	output, err := os.Create(filename)
	if err != nil {
		alert("Error while creating", filename, "-", err)
		return nil, err
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return nil, err
	}

	trace(n, "bytes downloaded.")
	return os.Open(filename)
}

func downloadFile(name string) (io.Reader, error) { // todo: need this here?

	u, err := url.Parse("https://" + name)
	if err != nil {
		panic(err)
	}

	local := support.Hash([]byte(name))

	if u.Host == "github.com" {
		bits := strings.Split(u.Path, "/")
		if len(bits) > 3 {
			nnn := make([]string, 0, 20)
			nnn = append(nnn, bits[:3]...)
			nnn = append(nnn, "raw")
			nnn = append(nnn, "master")
			nnn = append(nnn, bits[3:]...)

			u.Path = strings.Join(nnn, "/")
		}
	}

	return downloadFromUrl(u.String(), local)
}
