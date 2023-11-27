// Copyright 2023 The Connect Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: connectrpc/conformance/v1/service.proto

package conformancev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ConformanceService_Unary_FullMethodName         = "/connectrpc.conformance.v1.ConformanceService/Unary"
	ConformanceService_ServerStream_FullMethodName  = "/connectrpc.conformance.v1.ConformanceService/ServerStream"
	ConformanceService_ClientStream_FullMethodName  = "/connectrpc.conformance.v1.ConformanceService/ClientStream"
	ConformanceService_BidiStream_FullMethodName    = "/connectrpc.conformance.v1.ConformanceService/BidiStream"
	ConformanceService_Unimplemented_FullMethodName = "/connectrpc.conformance.v1.ConformanceService/Unimplemented"
)

// ConformanceServiceClient is the client API for ConformanceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConformanceServiceClient interface {
	// A unary operation. The request indicates the response headers and trailers
	// and also indicates either a response message or an error to send back.
	//
	// Response message data is specified as bytes. The service should echo back
	// request properties in the ConformancePayload and then include the message
	// data in the data field.
	//
	// Servers should allow the response definition to be unset in the request and
	// if it is, set no response headers or trailers and send back an empty response.
	Unary(ctx context.Context, in *UnaryRequest, opts ...grpc.CallOption) (*UnaryResponse, error)
	// A server-streaming operation. The request indicates the response headers,
	// response messages, trailers, and an optional error to send back. The
	// response data should be sent in the order indicated, and the server should
	// wait between sending response messages as indicated.
	//
	// Response message data is specified as bytes. The service should echo back
	// request properties in the first ConformancePayload, and then include the
	// message data in the data field. Subsequent messages after the first one
	// should contain only the data field.
	//
	// Servers should allow the response definition to be unset in the request and
	// if so, all responses should contain no response headers or trailers and
	// contain empty response data.
	ServerStream(ctx context.Context, in *ServerStreamRequest, opts ...grpc.CallOption) (ConformanceService_ServerStreamClient, error)
	// Servers should allow the response definition to be unset in the request and
	// if it is, set no response headers or trailers and send back empty response data.
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (ConformanceService_ClientStreamClient, error)
	// A bidirectional-streaming operation. The first request indicates the response
	// headers, response messages, trailers, and an optional error to send back.
	// The response data should be sent in the order indicated, and the server
	// should wait between sending response messages as indicated. If the
	// full_duplex field is true, the handler should read one request
	// and then send back one response, and then alternate, reading another
	// request and then sending back another response, etc. If the response_delay_ms
	// duration is specified, the server should wait that long in between sending each
	// response message. If both are specified, the server should wait the given
	// duration after reading the request before sending the corresponding
	// response.
	//
	// Response message data is specified as bytes and should be included in the
	// data field of the ConformancePayload in each response.
	//
	// If the full_duplex field is true, the service should echo back all request
	// properties in the first response including the last received request.
	// Subsequent responses should only echo back the last received request.
	//
	// If the full_duplex field is false, the service should echo back all request
	// properties, including all request messages in the order they were
	// received, in the ConformancePayload. Subsequent responses should only include
	// the message data in the data field.
	//
	// If the input stream is empty, the server should send a single response
	// message that includes no data and only the request properties (headers,
	// timeout).
	BidiStream(ctx context.Context, opts ...grpc.CallOption) (ConformanceService_BidiStreamClient, error)
	// A unary endpoint that the server should not implement and should instead
	// return an unimplemented error when invoked.
	Unimplemented(ctx context.Context, in *UnimplementedRequest, opts ...grpc.CallOption) (*UnimplementedResponse, error)
}

type conformanceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConformanceServiceClient(cc grpc.ClientConnInterface) ConformanceServiceClient {
	return &conformanceServiceClient{cc}
}

func (c *conformanceServiceClient) Unary(ctx context.Context, in *UnaryRequest, opts ...grpc.CallOption) (*UnaryResponse, error) {
	out := new(UnaryResponse)
	err := c.cc.Invoke(ctx, ConformanceService_Unary_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *conformanceServiceClient) ServerStream(ctx context.Context, in *ServerStreamRequest, opts ...grpc.CallOption) (ConformanceService_ServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConformanceService_ServiceDesc.Streams[0], ConformanceService_ServerStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &conformanceServiceServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConformanceService_ServerStreamClient interface {
	Recv() (*ServerStreamResponse, error)
	grpc.ClientStream
}

type conformanceServiceServerStreamClient struct {
	grpc.ClientStream
}

func (x *conformanceServiceServerStreamClient) Recv() (*ServerStreamResponse, error) {
	m := new(ServerStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *conformanceServiceClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (ConformanceService_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConformanceService_ServiceDesc.Streams[1], ConformanceService_ClientStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &conformanceServiceClientStreamClient{stream}
	return x, nil
}

type ConformanceService_ClientStreamClient interface {
	Send(*ClientStreamRequest) error
	CloseAndRecv() (*ClientStreamResponse, error)
	grpc.ClientStream
}

type conformanceServiceClientStreamClient struct {
	grpc.ClientStream
}

func (x *conformanceServiceClientStreamClient) Send(m *ClientStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *conformanceServiceClientStreamClient) CloseAndRecv() (*ClientStreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ClientStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *conformanceServiceClient) BidiStream(ctx context.Context, opts ...grpc.CallOption) (ConformanceService_BidiStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConformanceService_ServiceDesc.Streams[2], ConformanceService_BidiStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &conformanceServiceBidiStreamClient{stream}
	return x, nil
}

type ConformanceService_BidiStreamClient interface {
	Send(*BidiStreamRequest) error
	Recv() (*BidiStreamResponse, error)
	grpc.ClientStream
}

type conformanceServiceBidiStreamClient struct {
	grpc.ClientStream
}

func (x *conformanceServiceBidiStreamClient) Send(m *BidiStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *conformanceServiceBidiStreamClient) Recv() (*BidiStreamResponse, error) {
	m := new(BidiStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *conformanceServiceClient) Unimplemented(ctx context.Context, in *UnimplementedRequest, opts ...grpc.CallOption) (*UnimplementedResponse, error) {
	out := new(UnimplementedResponse)
	err := c.cc.Invoke(ctx, ConformanceService_Unimplemented_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConformanceServiceServer is the server API for ConformanceService service.
// All implementations must embed UnimplementedConformanceServiceServer
// for forward compatibility
type ConformanceServiceServer interface {
	// A unary operation. The request indicates the response headers and trailers
	// and also indicates either a response message or an error to send back.
	//
	// Response message data is specified as bytes. The service should echo back
	// request properties in the ConformancePayload and then include the message
	// data in the data field.
	//
	// Servers should allow the response definition to be unset in the request and
	// if it is, set no response headers or trailers and send back an empty response.
	Unary(context.Context, *UnaryRequest) (*UnaryResponse, error)
	// A server-streaming operation. The request indicates the response headers,
	// response messages, trailers, and an optional error to send back. The
	// response data should be sent in the order indicated, and the server should
	// wait between sending response messages as indicated.
	//
	// Response message data is specified as bytes. The service should echo back
	// request properties in the first ConformancePayload, and then include the
	// message data in the data field. Subsequent messages after the first one
	// should contain only the data field.
	//
	// Servers should allow the response definition to be unset in the request and
	// if so, all responses should contain no response headers or trailers and
	// contain empty response data.
	ServerStream(*ServerStreamRequest, ConformanceService_ServerStreamServer) error
	// Servers should allow the response definition to be unset in the request and
	// if it is, set no response headers or trailers and send back empty response data.
	ClientStream(ConformanceService_ClientStreamServer) error
	// A bidirectional-streaming operation. The first request indicates the response
	// headers, response messages, trailers, and an optional error to send back.
	// The response data should be sent in the order indicated, and the server
	// should wait between sending response messages as indicated. If the
	// full_duplex field is true, the handler should read one request
	// and then send back one response, and then alternate, reading another
	// request and then sending back another response, etc. If the response_delay_ms
	// duration is specified, the server should wait that long in between sending each
	// response message. If both are specified, the server should wait the given
	// duration after reading the request before sending the corresponding
	// response.
	//
	// Response message data is specified as bytes and should be included in the
	// data field of the ConformancePayload in each response.
	//
	// If the full_duplex field is true, the service should echo back all request
	// properties in the first response including the last received request.
	// Subsequent responses should only echo back the last received request.
	//
	// If the full_duplex field is false, the service should echo back all request
	// properties, including all request messages in the order they were
	// received, in the ConformancePayload. Subsequent responses should only include
	// the message data in the data field.
	//
	// If the input stream is empty, the server should send a single response
	// message that includes no data and only the request properties (headers,
	// timeout).
	BidiStream(ConformanceService_BidiStreamServer) error
	// A unary endpoint that the server should not implement and should instead
	// return an unimplemented error when invoked.
	Unimplemented(context.Context, *UnimplementedRequest) (*UnimplementedResponse, error)
	mustEmbedUnimplementedConformanceServiceServer()
}

// UnimplementedConformanceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConformanceServiceServer struct {
}

func (UnimplementedConformanceServiceServer) Unary(context.Context, *UnaryRequest) (*UnaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unary not implemented")
}
func (UnimplementedConformanceServiceServer) ServerStream(*ServerStreamRequest, ConformanceService_ServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStream not implemented")
}
func (UnimplementedConformanceServiceServer) ClientStream(ConformanceService_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (UnimplementedConformanceServiceServer) BidiStream(ConformanceService_BidiStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method BidiStream not implemented")
}
func (UnimplementedConformanceServiceServer) Unimplemented(context.Context, *UnimplementedRequest) (*UnimplementedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unimplemented not implemented")
}
func (UnimplementedConformanceServiceServer) mustEmbedUnimplementedConformanceServiceServer() {}

// UnsafeConformanceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConformanceServiceServer will
// result in compilation errors.
type UnsafeConformanceServiceServer interface {
	mustEmbedUnimplementedConformanceServiceServer()
}

func RegisterConformanceServiceServer(s grpc.ServiceRegistrar, srv ConformanceServiceServer) {
	s.RegisterService(&ConformanceService_ServiceDesc, srv)
}

func _ConformanceService_Unary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConformanceServiceServer).Unary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConformanceService_Unary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConformanceServiceServer).Unary(ctx, req.(*UnaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConformanceService_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ServerStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConformanceServiceServer).ServerStream(m, &conformanceServiceServerStreamServer{stream})
}

type ConformanceService_ServerStreamServer interface {
	Send(*ServerStreamResponse) error
	grpc.ServerStream
}

type conformanceServiceServerStreamServer struct {
	grpc.ServerStream
}

func (x *conformanceServiceServerStreamServer) Send(m *ServerStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ConformanceService_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ConformanceServiceServer).ClientStream(&conformanceServiceClientStreamServer{stream})
}

type ConformanceService_ClientStreamServer interface {
	SendAndClose(*ClientStreamResponse) error
	Recv() (*ClientStreamRequest, error)
	grpc.ServerStream
}

type conformanceServiceClientStreamServer struct {
	grpc.ServerStream
}

func (x *conformanceServiceClientStreamServer) SendAndClose(m *ClientStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *conformanceServiceClientStreamServer) Recv() (*ClientStreamRequest, error) {
	m := new(ClientStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ConformanceService_BidiStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ConformanceServiceServer).BidiStream(&conformanceServiceBidiStreamServer{stream})
}

type ConformanceService_BidiStreamServer interface {
	Send(*BidiStreamResponse) error
	Recv() (*BidiStreamRequest, error)
	grpc.ServerStream
}

type conformanceServiceBidiStreamServer struct {
	grpc.ServerStream
}

func (x *conformanceServiceBidiStreamServer) Send(m *BidiStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *conformanceServiceBidiStreamServer) Recv() (*BidiStreamRequest, error) {
	m := new(BidiStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ConformanceService_Unimplemented_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnimplementedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConformanceServiceServer).Unimplemented(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConformanceService_Unimplemented_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConformanceServiceServer).Unimplemented(ctx, req.(*UnimplementedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConformanceService_ServiceDesc is the grpc.ServiceDesc for ConformanceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConformanceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "connectrpc.conformance.v1.ConformanceService",
	HandlerType: (*ConformanceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Unary",
			Handler:    _ConformanceService_Unary_Handler,
		},
		{
			MethodName: "Unimplemented",
			Handler:    _ConformanceService_Unimplemented_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStream",
			Handler:       _ConformanceService_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _ConformanceService_ClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "BidiStream",
			Handler:       _ConformanceService_BidiStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "connectrpc/conformance/v1/service.proto",
}
