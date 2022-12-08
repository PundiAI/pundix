// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/ibc/applications/transfer/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	types "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("fx/ibc/applications/transfer/v1/query.proto", fileDescriptor_569f08cc402420ba)
}

var fileDescriptor_569f08cc402420ba = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x3f, 0x4b, 0x03, 0x31,
	0x1c, 0x86, 0xef, 0x40, 0x0b, 0x46, 0x5c, 0x32, 0x76, 0x88, 0xe0, 0x24, 0x08, 0x89, 0xb5, 0xf5,
	0x0f, 0x6e, 0x8a, 0x82, 0xa3, 0x4a, 0x27, 0xb7, 0x5c, 0x9a, 0xda, 0x80, 0x77, 0x49, 0xf3, 0xcb,
	0xd5, 0x2b, 0x82, 0x93, 0x1f, 0xc0, 0x8f, 0xe5, 0xd8, 0xd1, 0x51, 0xee, 0xbe, 0x83, 0xb3, 0xf4,
	0x24, 0x57, 0x5d, 0xca, 0x65, 0x3a, 0x2e, 0xbc, 0xcf, 0xfb, 0xbc, 0x43, 0x82, 0x0e, 0xc6, 0x05,
	0x53, 0x89, 0x60, 0xdc, 0x98, 0x27, 0x25, 0xb8, 0x53, 0x3a, 0x03, 0xe6, 0x2c, 0xcf, 0x60, 0x2c,
	0x2d, 0x9b, 0xf5, 0xd8, 0x34, 0x97, 0x76, 0x4e, 0x8d, 0xd5, 0x4e, 0xe3, 0xdd, 0x71, 0x41, 0x55,
	0x22, 0xe8, 0xdf, 0x30, 0xf5, 0x61, 0x3a, 0xeb, 0x75, 0xf7, 0xdb, 0x56, 0x1d, 0x7d, 0x6f, 0xa0,
	0xcd, 0xbb, 0xe5, 0x3f, 0x7e, 0x41, 0xe8, 0x4a, 0x66, 0x3a, 0x1d, 0x5a, 0x2e, 0x24, 0x1e, 0xac,
	0x15, 0xd0, 0x1a, 0x59, 0xc5, 0xef, 0xe5, 0x34, 0x97, 0xe0, 0xba, 0xc7, 0x81, 0x14, 0x18, 0x9d,
	0x81, 0xdc, 0x8b, 0xf0, 0x2b, 0xda, 0x5e, 0x9d, 0x03, 0x0e, 0xeb, 0x01, 0xaf, 0x3f, 0x09, 0xc5,
	0x1a, 0xbf, 0x46, 0x9d, 0x5b, 0x6e, 0x79, 0x0a, 0xf8, 0xb0, 0x45, 0xc7, 0x6f, 0xd4, 0x5b, 0x7b,
	0x01, 0x44, 0x23, 0x2c, 0xd0, 0x56, 0xbd, 0xe4, 0x86, 0xc3, 0x04, 0xf7, 0xdb, 0xee, 0x5e, 0xa6,
	0xbd, 0x76, 0x10, 0x06, 0x35, 0xe6, 0xb7, 0x18, 0xed, 0x5c, 0x83, 0xb0, 0xfa, 0xf9, 0x62, 0x34,
	0xb2, 0x12, 0x00, 0x9f, 0xb6, 0x68, 0xfa, 0x47, 0xf8, 0x09, 0x67, 0xe1, 0xa0, 0x9f, 0x71, 0x39,
	0xfc, 0x28, 0x49, 0xbc, 0x28, 0x49, 0xfc, 0x55, 0x92, 0xf8, 0xbd, 0x22, 0xd1, 0xa2, 0x22, 0xd1,
	0x67, 0x45, 0xa2, 0x87, 0xf3, 0x47, 0xe5, 0x26, 0x79, 0x42, 0x85, 0x4e, 0x99, 0xc9, 0xb3, 0x91,
	0x2a, 0xfc, 0x67, 0xdd, 0x13, 0x71, 0x73, 0x23, 0x21, 0xe9, 0xd4, 0xb7, 0xba, 0xff, 0x13, 0x00,
	0x00, 0xff, 0xff, 0xa8, 0xfd, 0x79, 0x14, 0x4f, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// DenomTrace queries a denomination trace information.
	DenomTrace(ctx context.Context, in *types.QueryDenomTraceRequest, opts ...grpc.CallOption) (*types.QueryDenomTraceResponse, error)
	// DenomTraces queries all denomination traces.
	DenomTraces(ctx context.Context, in *types.QueryDenomTracesRequest, opts ...grpc.CallOption) (*types.QueryDenomTracesResponse, error)
	// Params queries all parameters of the ibc-transfer module.
	Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error)
	// DenomHash queries a denomination hash information.
	DenomHash(ctx context.Context, in *types.QueryDenomHashRequest, opts ...grpc.CallOption) (*types.QueryDenomHashResponse, error)
	// EscrowAddress returns the escrow address for a particular port and channel
	// id.
	EscrowAddress(ctx context.Context, in *types.QueryEscrowAddressRequest, opts ...grpc.CallOption) (*types.QueryEscrowAddressResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) DenomTrace(ctx context.Context, in *types.QueryDenomTraceRequest, opts ...grpc.CallOption) (*types.QueryDenomTraceResponse, error) {
	out := new(types.QueryDenomTraceResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/DenomTrace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomTraces(ctx context.Context, in *types.QueryDenomTracesRequest, opts ...grpc.CallOption) (*types.QueryDenomTracesResponse, error) {
	out := new(types.QueryDenomTracesResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/DenomTraces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error) {
	out := new(types.QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomHash(ctx context.Context, in *types.QueryDenomHashRequest, opts ...grpc.CallOption) (*types.QueryDenomHashResponse, error) {
	out := new(types.QueryDenomHashResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/DenomHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EscrowAddress(ctx context.Context, in *types.QueryEscrowAddressRequest, opts ...grpc.CallOption) (*types.QueryEscrowAddressResponse, error) {
	out := new(types.QueryEscrowAddressResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/EscrowAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// DenomTrace queries a denomination trace information.
	DenomTrace(context.Context, *types.QueryDenomTraceRequest) (*types.QueryDenomTraceResponse, error)
	// DenomTraces queries all denomination traces.
	DenomTraces(context.Context, *types.QueryDenomTracesRequest) (*types.QueryDenomTracesResponse, error)
	// Params queries all parameters of the ibc-transfer module.
	Params(context.Context, *types.QueryParamsRequest) (*types.QueryParamsResponse, error)
	// DenomHash queries a denomination hash information.
	DenomHash(context.Context, *types.QueryDenomHashRequest) (*types.QueryDenomHashResponse, error)
	// EscrowAddress returns the escrow address for a particular port and channel
	// id.
	EscrowAddress(context.Context, *types.QueryEscrowAddressRequest) (*types.QueryEscrowAddressResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) DenomTrace(ctx context.Context, req *types.QueryDenomTraceRequest) (*types.QueryDenomTraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomTrace not implemented")
}
func (*UnimplementedQueryServer) DenomTraces(ctx context.Context, req *types.QueryDenomTracesRequest) (*types.QueryDenomTracesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomTraces not implemented")
}
func (*UnimplementedQueryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) DenomHash(ctx context.Context, req *types.QueryDenomHashRequest) (*types.QueryDenomHashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomHash not implemented")
}
func (*UnimplementedQueryServer) EscrowAddress(ctx context.Context, req *types.QueryEscrowAddressRequest) (*types.QueryEscrowAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EscrowAddress not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_DenomTrace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryDenomTraceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomTrace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/DenomTrace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomTrace(ctx, req.(*types.QueryDenomTraceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomTraces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryDenomTracesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomTraces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/DenomTraces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomTraces(ctx, req.(*types.QueryDenomTracesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*types.QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryDenomHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/DenomHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomHash(ctx, req.(*types.QueryDenomHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EscrowAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryEscrowAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EscrowAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/EscrowAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EscrowAddress(ctx, req.(*types.QueryEscrowAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fx.ibc.applications.transfer.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DenomTrace",
			Handler:    _Query_DenomTrace_Handler,
		},
		{
			MethodName: "DenomTraces",
			Handler:    _Query_DenomTraces_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "DenomHash",
			Handler:    _Query_DenomHash_Handler,
		},
		{
			MethodName: "EscrowAddress",
			Handler:    _Query_EscrowAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fx/ibc/applications/transfer/v1/query.proto",
}
