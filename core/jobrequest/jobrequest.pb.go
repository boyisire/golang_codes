// Code generated by protoc-gen-go. DO NOT EDIT.
// source: jobrequest.proto

/*
Package jobrequest is a generated protocol buffer package.

It is generated from these files:
	jobrequest.proto

It has these top-level messages:
	JobRequest
*/
package jobrequest

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type JobRequest struct {
	Command string            `protobuf:"bytes,1,opt,name=command" json:"command,omitempty"`
	Data    map[string]string `protobuf:"bytes,4,rep,name=data" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *JobRequest) Reset()                    { *m = JobRequest{} }
func (m *JobRequest) String() string            { return proto.CompactTextString(m) }
func (*JobRequest) ProtoMessage()               {}
func (*JobRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *JobRequest) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *JobRequest) GetData() map[string]string {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*JobRequest)(nil), "JobRequest")
}

func init() { proto.RegisterFile("jobrequest.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 145 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0xca, 0x4f, 0x2a,
	0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xea, 0x62, 0xe4,
	0xe2, 0xf2, 0xca, 0x4f, 0x0a, 0x82, 0x08, 0x0a, 0x49, 0x70, 0xb1, 0x27, 0xe7, 0xe7, 0xe6, 0x26,
	0xe6, 0xa5, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x42, 0x9a, 0x5c, 0x2c, 0x29,
	0x89, 0x25, 0x89, 0x12, 0x2c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0xa2, 0x7a, 0x08, 0x4d, 0x7a, 0x2e,
	0x89, 0x25, 0x89, 0xae, 0x79, 0x25, 0x45, 0x95, 0x41, 0x60, 0x25, 0x52, 0xe6, 0x5c, 0x9c, 0x70,
	0x21, 0x21, 0x01, 0x2e, 0xe6, 0xec, 0xd4, 0x4a, 0xa8, 0x69, 0x20, 0xa6, 0x90, 0x08, 0x17, 0x6b,
	0x59, 0x62, 0x4e, 0x69, 0xaa, 0x04, 0x13, 0x58, 0x0c, 0xc2, 0xb1, 0x62, 0xb2, 0x60, 0x4c, 0x62,
	0x03, 0xbb, 0xc9, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xee, 0x7e, 0x49, 0x76, 0xa7, 0x00, 0x00,
	0x00,
}