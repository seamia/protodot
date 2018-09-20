// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type PBS struct {
	Package    string
	Protoname  string
	AppVersion string
	Timestamp  string
	Selection  string
	Options    string
}

type Relationship struct {
	From  string
	Field string

	To     UniqueName
	ToName string
	ToType FullName
}

type Cluster struct {
	ProtoNameKosher string
	ProtoName       string
	ShortName       string
}

type EnumPayload struct {
	Name     string
	Value    string
	Unique   UniqueName
	FullName FullName
}

type RPC struct {
	Name           string
	RequestType    string
	ReturnsType    string
	StreamsRequest string
	StreamsReturns string
}

type ServicePayload struct {
	Name     string
	Unique   UniqueName
	FullName FullName
}
