// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/emicklei/proto"
)

//----------------------------------------------------------------------------------------------------------------------
// additions/extensions to "github.com/emicklei/proto"

func WithPackage(apply func(*proto.Package)) proto.Handler {
	return func(v proto.Visitee) {
		if s, ok := v.(*proto.Package); ok {
			apply(s)
		}
	}
}

func WithImport(apply func(*proto.Import)) proto.Handler {
	return func(v proto.Visitee) {
		if s, ok := v.(*proto.Import); ok {
			apply(s)
		}
	}
}

func WithSyntax(apply func(*proto.Syntax)) proto.Handler {
	return func(v proto.Visitee) {
		if s, ok := v.(*proto.Syntax); ok {
			apply(s)
		}
	}
}
