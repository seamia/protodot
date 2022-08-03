// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/emicklei/proto"
	"github.com/seamia/protodot/plus"
)

//----------------------------------------------------------------------------------------------------------------------
// presentation
type table struct {
	buffer bytes.Buffer
	name   string
}

func newTable(name string, full FullName, unique UniqueName, style string) *table {
	t := table{
		name: name,
	}

	entry := OneOfEntry{
		Name:   name,
		Unique: unique,
		Type:   string(full),
	}
	if err := plus.ApplyTemplate("message.prefix", &t, entry); err != nil {
		alert("failed to render", err)
	}

	return &t
}

func (t *table) Write(p []byte) (n int, err error) {
	return t.buffer.Write(p)
}

func (t *table) WriteString(s string) {
	t.Write([]byte(s))
}

var kind2entry = map[Kind]string{
	Simple:  "entry.simple",
	Message: "entry.message",
	Enum:    "entry.enum",
	Missing: "entry.missing",
	Comment: "entry.comment",
}

func (t *table) addRow(repeated, typ, name, ordinal string, kind Kind, extra string) {
	tmplName := kind2entry[kind]
	if len(tmplName) > 0 {
		entry := OneOfEntry{
			Name:    name,
			Type:    typ,
			Ordinal: ordinal,
			Prefix:  repeated,
			Extra:   extra,
		}
		if err := plus.ApplyTemplate(tmplName, t, entry); err != nil {
			alert("failed to render", err)
		}
	} else {
		alert("unhandled kind", kind)
	}
}

var kind2map = map[Kind]string{
	Simple:  "map.simple",
	Message: "map.message",
	Enum:    "map.enum",
	Missing: "map.missing",
}

func (t *table) addMapRow(name, keyType, typ, ordinal string, kind Kind) {

	tmplName := kind2map[kind]
	if len(tmplName) > 0 {
		entry := OneOfEntry{
			Name:    name,
			Type:    typ,
			Ordinal: ordinal,
			Prefix:  "",
			KeyType: keyType,
		}
		if err := plus.ApplyTemplate(tmplName, t, entry); err != nil {
			alert("failed to render", err)
		}
	} else {
		alert("unhandled kind:", kind)
	}
}

type OneOfEntry struct {
	Unique  UniqueName
	Name    string
	KeyType string
	Type    string
	Ordinal string
	Prefix  string
	Extra   string
}

var kind2template = map[Kind]string{
	Simple:  "oneof.entry.simple",
	Message: "oneof.entry.message",
	Enum:    "oneof.entry.enum",
	Missing: "oneof.entry.missing",
}

func (t *table) addOneof(fullname FullName, what *proto.Oneof, pbs *pbstate) {

	entry := OneOfEntry{
		Name: what.Name,
	}
	if err := plus.ApplyTemplate("oneof.entry.prefix", t, entry); err != nil {
		alert("failed to render", err)
	}

	for _, element := range what.Elements {
		switch actual := element.(type) {
		case *proto.OneOfField:
			if tmplName := kind2template[pbs.getKind(fullname, OriginalName(actual.Type))]; len(tmplName) > 0 {
				payload := OneOfEntry{
					Name:    actual.Name,
					Type:    actual.Type,
					Ordinal: strconv.Itoa(actual.Sequence),
				}
				if err := plus.ApplyTemplate(tmplName, t, payload); err != nil {
					alert("failed to render", err)
				}
			} else {
				alert("failed to get template name")
			}

		case *proto.Option:
			ignoring("ignoring options for now")

		case *proto.Comment:
			ignoring("ignoring comments for now")

		case *proto.Group:
			ignoring("ignoring group for now")

		default:
			rname := reflect.TypeOf(actual).Elem().Name()
			unhandled("\t", "UNKNOWN1", actual, "", rname)
		}
	}

	if err := plus.ApplyTemplate("oneof.entry.suffix", t, entry); err != nil {
		alert("failed to render", err)
	}
}

func (t *table) generate() string {

	entry := OneOfEntry{Name: t.name}
	if err := plus.ApplyTemplate("message.suffix", t, entry); err != nil {
		alert("failed to render", err)
	}

	return t.buffer.String()
}
