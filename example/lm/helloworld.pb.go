// Code generated by protoc-gen-go.
// source: helloworld.proto
// DO NOT EDIT!

/*
Package lm is a generated protocol buffer package.

It is generated from these files:
	helloworld.proto

It has these top-level messages:
	Helloworld
*/
package lm

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Helloworld struct {
	Id               *int32  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Str              *string `protobuf:"bytes,2,req,name=str" json:"str,omitempty"`
	Opt              *int32  `protobuf:"varint,3,opt,name=opt" json:"opt,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Helloworld) Reset()         { *m = Helloworld{} }
func (m *Helloworld) String() string { return proto.CompactTextString(m) }
func (*Helloworld) ProtoMessage()    {}

func (m *Helloworld) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Helloworld) GetStr() string {
	if m != nil && m.Str != nil {
		return *m.Str
	}
	return ""
}

func (m *Helloworld) GetOpt() int32 {
	if m != nil && m.Opt != nil {
		return *m.Opt
	}
	return 0
}

func init() {
}