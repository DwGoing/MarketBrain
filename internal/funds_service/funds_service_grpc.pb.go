// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: proto/funds_service.proto

package funds_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FundsServiceClient is the client API for FundsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FundsServiceClient interface {
	GetCollectionWallet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetCollectionWalletResponse, error)
	GetRechargeWallet(ctx context.Context, in *GetRechargeWalletRequest, opts ...grpc.CallOption) (*GetRechargeWalletResponse, error)
	GetRechargeRecords(ctx context.Context, in *GetRechargeRecordsRequest, opts ...grpc.CallOption) (*GetRechargeRecordsResponse, error)
	FundsCollect(ctx context.Context, in *FundsCollectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type fundsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFundsServiceClient(cc grpc.ClientConnInterface) FundsServiceClient {
	return &fundsServiceClient{cc}
}

func (c *fundsServiceClient) GetCollectionWallet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetCollectionWalletResponse, error) {
	out := new(GetCollectionWalletResponse)
	err := c.cc.Invoke(ctx, "/FundsService/GetCollectionWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fundsServiceClient) GetRechargeWallet(ctx context.Context, in *GetRechargeWalletRequest, opts ...grpc.CallOption) (*GetRechargeWalletResponse, error) {
	out := new(GetRechargeWalletResponse)
	err := c.cc.Invoke(ctx, "/FundsService/GetRechargeWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fundsServiceClient) GetRechargeRecords(ctx context.Context, in *GetRechargeRecordsRequest, opts ...grpc.CallOption) (*GetRechargeRecordsResponse, error) {
	out := new(GetRechargeRecordsResponse)
	err := c.cc.Invoke(ctx, "/FundsService/GetRechargeRecords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fundsServiceClient) FundsCollect(ctx context.Context, in *FundsCollectRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/FundsService/FundsCollect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FundsServiceServer is the server API for FundsService service.
// All implementations must embed UnimplementedFundsServiceServer
// for forward compatibility
type FundsServiceServer interface {
	GetCollectionWallet(context.Context, *emptypb.Empty) (*GetCollectionWalletResponse, error)
	GetRechargeWallet(context.Context, *GetRechargeWalletRequest) (*GetRechargeWalletResponse, error)
	GetRechargeRecords(context.Context, *GetRechargeRecordsRequest) (*GetRechargeRecordsResponse, error)
	FundsCollect(context.Context, *FundsCollectRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedFundsServiceServer()
}

// UnimplementedFundsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFundsServiceServer struct {
}

func (UnimplementedFundsServiceServer) GetCollectionWallet(context.Context, *emptypb.Empty) (*GetCollectionWalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollectionWallet not implemented")
}
func (UnimplementedFundsServiceServer) GetRechargeWallet(context.Context, *GetRechargeWalletRequest) (*GetRechargeWalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRechargeWallet not implemented")
}
func (UnimplementedFundsServiceServer) GetRechargeRecords(context.Context, *GetRechargeRecordsRequest) (*GetRechargeRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRechargeRecords not implemented")
}
func (UnimplementedFundsServiceServer) FundsCollect(context.Context, *FundsCollectRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FundsCollect not implemented")
}
func (UnimplementedFundsServiceServer) mustEmbedUnimplementedFundsServiceServer() {}

// UnsafeFundsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FundsServiceServer will
// result in compilation errors.
type UnsafeFundsServiceServer interface {
	mustEmbedUnimplementedFundsServiceServer()
}

func RegisterFundsServiceServer(s grpc.ServiceRegistrar, srv FundsServiceServer) {
	s.RegisterService(&FundsService_ServiceDesc, srv)
}

func _FundsService_GetCollectionWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FundsServiceServer).GetCollectionWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FundsService/GetCollectionWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FundsServiceServer).GetCollectionWallet(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _FundsService_GetRechargeWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRechargeWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FundsServiceServer).GetRechargeWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FundsService/GetRechargeWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FundsServiceServer).GetRechargeWallet(ctx, req.(*GetRechargeWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FundsService_GetRechargeRecords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRechargeRecordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FundsServiceServer).GetRechargeRecords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FundsService/GetRechargeRecords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FundsServiceServer).GetRechargeRecords(ctx, req.(*GetRechargeRecordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FundsService_FundsCollect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FundsCollectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FundsServiceServer).FundsCollect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FundsService/FundsCollect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FundsServiceServer).FundsCollect(ctx, req.(*FundsCollectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FundsService_ServiceDesc is the grpc.ServiceDesc for FundsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FundsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FundsService",
	HandlerType: (*FundsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCollectionWallet",
			Handler:    _FundsService_GetCollectionWallet_Handler,
		},
		{
			MethodName: "GetRechargeWallet",
			Handler:    _FundsService_GetRechargeWallet_Handler,
		},
		{
			MethodName: "GetRechargeRecords",
			Handler:    _FundsService_GetRechargeRecords_Handler,
		},
		{
			MethodName: "FundsCollect",
			Handler:    _FundsService_FundsCollect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/funds_service.proto",
}