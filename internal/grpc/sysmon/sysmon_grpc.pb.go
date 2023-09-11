// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: sysmon.proto

package sysmon

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

// SysmonDataClient is the client API for SysmonData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SysmonDataClient interface {
	SendSysmonDataToClient(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (SysmonData_SendSysmonDataToClientClient, error)
}

type sysmonDataClient struct {
	cc grpc.ClientConnInterface
}

func NewSysmonDataClient(cc grpc.ClientConnInterface) SysmonDataClient {
	return &sysmonDataClient{cc}
}

func (c *sysmonDataClient) SendSysmonDataToClient(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (SysmonData_SendSysmonDataToClientClient, error) {
	stream, err := c.cc.NewStream(ctx, &SysmonData_ServiceDesc.Streams[0], "/sysmon.SysmonData/SendSysmonDataToClient", opts...)
	if err != nil {
		return nil, err
	}
	x := &sysmonDataSendSysmonDataToClientClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SysmonData_SendSysmonDataToClientClient interface {
	Recv() (*Data, error)
	grpc.ClientStream
}

type sysmonDataSendSysmonDataToClientClient struct {
	grpc.ClientStream
}

func (x *sysmonDataSendSysmonDataToClientClient) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SysmonDataServer is the server API for SysmonData service.
// All implementations should embed UnimplementedSysmonDataServer
// for forward compatibility
type SysmonDataServer interface {
	SendSysmonDataToClient(*QueryParam, SysmonData_SendSysmonDataToClientServer) error
}

// UnimplementedSysmonDataServer should be embedded to have forward compatible implementations.
type UnimplementedSysmonDataServer struct {
}

func (UnimplementedSysmonDataServer) SendSysmonDataToClient(*QueryParam, SysmonData_SendSysmonDataToClientServer) error {
	return status.Errorf(codes.Unimplemented, "method SendSysmonDataToClient not implemented")
}

// UnsafeSysmonDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SysmonDataServer will
// result in compilation errors.
type UnsafeSysmonDataServer interface {
	mustEmbedUnimplementedSysmonDataServer()
}

func RegisterSysmonDataServer(s grpc.ServiceRegistrar, srv SysmonDataServer) {
	s.RegisterService(&SysmonData_ServiceDesc, srv)
}

func _SysmonData_SendSysmonDataToClient_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryParam)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SysmonDataServer).SendSysmonDataToClient(m, &sysmonDataSendSysmonDataToClientServer{stream})
}

type SysmonData_SendSysmonDataToClientServer interface {
	Send(*Data) error
	grpc.ServerStream
}

type sysmonDataSendSysmonDataToClientServer struct {
	grpc.ServerStream
}

func (x *sysmonDataSendSysmonDataToClientServer) Send(m *Data) error {
	return x.ServerStream.SendMsg(m)
}

// SysmonData_ServiceDesc is the grpc.ServiceDesc for SysmonData service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SysmonData_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sysmon.SysmonData",
	HandlerType: (*SysmonDataServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendSysmonDataToClient",
			Handler:       _SysmonData_SendSysmonDataToClient_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "sysmon.proto",
}