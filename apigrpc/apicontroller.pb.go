// Code generated by protoc-gen-go. DO NOT EDIT.
// source: apicontroller.proto

/*
Package apigrpc is a generated protocol buffer package.

It is generated from these files:
	apicontroller.proto

It has these top-level messages:
	Request
	Reply
*/
package apigrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type Request struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Version string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
	Dc      string `protobuf:"bytes,4,opt,name=dc" json:"dc,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Request) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Request) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Request) GetDc() string {
	if m != nil {
		return m.Dc
	}
	return ""
}

// The response message containing the greetings
type Reply struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Data string `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *Reply) Reset()                    { *m = Reply{} }
func (m *Reply) String() string            { return proto.CompactTextString(m) }
func (*Reply) ProtoMessage()               {}
func (*Reply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Reply) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Reply) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "apigrpc.Request")
	proto.RegisterType((*Reply)(nil), "apigrpc.Reply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SendHelp service

type SendHelpClient interface {
	// Sends a greeting
	SendHelp(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
}

type sendHelpClient struct {
	cc *grpc.ClientConn
}

func NewSendHelpClient(cc *grpc.ClientConn) SendHelpClient {
	return &sendHelpClient{cc}
}

func (c *sendHelpClient) SendHelp(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/apigrpc.SendHelp/SendHelp", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SendHelp service

type SendHelpServer interface {
	// Sends a greeting
	SendHelp(context.Context, *Request) (*Reply, error)
}

func RegisterSendHelpServer(s *grpc.Server, srv SendHelpServer) {
	s.RegisterService(&_SendHelp_serviceDesc, srv)
}

func _SendHelp_SendHelp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendHelpServer).SendHelp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apigrpc.SendHelp/SendHelp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendHelpServer).SendHelp(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _SendHelp_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apigrpc.SendHelp",
	HandlerType: (*SendHelpServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendHelp",
			Handler:    _SendHelp_SendHelp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apicontroller.proto",
}

func init() { proto.RegisterFile("apicontroller.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xb1, 0x6b, 0xc4, 0x20,
	0x18, 0xc5, 0x9b, 0x34, 0x6d, 0xda, 0x6f, 0x08, 0xc5, 0x2e, 0xd2, 0xa9, 0x64, 0x2a, 0xb4, 0x38,
	0xb4, 0x5b, 0xc7, 0x40, 0xa1, 0x63, 0xc8, 0x0d, 0x37, 0x7b, 0xfa, 0x71, 0x08, 0x46, 0x3d, 0x63,
	0x0e, 0xf2, 0xdf, 0x1f, 0x7a, 0x09, 0x24, 0x70, 0xdb, 0xf3, 0xa7, 0xef, 0xf1, 0x7c, 0xf0, 0xca,
	0x9d, 0x12, 0xd6, 0x04, 0x6f, 0xb5, 0x46, 0xcf, 0x9c, 0xb7, 0xc1, 0x92, 0x92, 0x3b, 0x75, 0xf4,
	0x4e, 0xd4, 0x7b, 0x28, 0x3b, 0x3c, 0x8d, 0x38, 0x04, 0x52, 0x41, 0xae, 0x24, 0xcd, 0xde, 0xb3,
	0x8f, 0xe7, 0x2e, 0x57, 0x92, 0x10, 0x28, 0x0c, 0xef, 0x91, 0xe6, 0x89, 0x24, 0x4d, 0x28, 0x94,
	0x67, 0xf4, 0x83, 0xb2, 0x86, 0xde, 0x27, 0xbc, 0x1c, 0xa3, 0x5b, 0x0a, 0x5a, 0x5c, 0xdd, 0x52,
	0xd4, 0x9f, 0xf0, 0xd0, 0xa1, 0xd3, 0xd3, 0xad, 0x58, 0xc9, 0x03, 0x5f, 0x62, 0xa3, 0xfe, 0xfe,
	0x85, 0xa7, 0x1d, 0x1a, 0xf9, 0x8f, 0xda, 0x11, 0xb6, 0xd2, 0x2f, 0x6c, 0xee, 0xc9, 0xe6, 0x92,
	0x6f, 0xd5, 0x8a, 0x38, 0x3d, 0xd5, 0x77, 0xcd, 0x1f, 0x7c, 0x09, 0xdb, 0x33, 0x6e, 0xa6, 0xd1,
	0xc4, 0x4b, 0xb6, 0xfa, 0x6c, 0x7a, 0xb8, 0x65, 0xcd, 0x76, 0x0f, 0xd9, 0xc6, 0x3d, 0xda, 0xec,
	0xf0, 0x98, 0x86, 0xf9, 0xb9, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd6, 0x93, 0xa0, 0x29, 0x2f, 0x01,
	0x00, 0x00,
}
