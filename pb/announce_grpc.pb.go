// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ServiceDiscoveryClient is the client API for ServiceDiscovery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceDiscoveryClient interface {
	// Announce that an app exists
	Announce(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
}

type serviceDiscoveryClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceDiscoveryClient(cc grpc.ClientConnInterface) ServiceDiscoveryClient {
	return &serviceDiscoveryClient{cc}
}

func (c *serviceDiscoveryClient) Announce(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/ServiceDiscovery/Announce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceDiscoveryServer is the server API for ServiceDiscovery service.
// All implementations must embed UnimplementedServiceDiscoveryServer
// for forward compatibility
type ServiceDiscoveryServer interface {
	// Announce that an app exists
	Announce(context.Context, *RegisterRequest) (*RegisterResponse, error)
	mustEmbedUnimplementedServiceDiscoveryServer()
}

// UnimplementedServiceDiscoveryServer must be embedded to have forward compatible implementations.
type UnimplementedServiceDiscoveryServer struct {
}

func (UnimplementedServiceDiscoveryServer) Announce(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Announce not implemented")
}
func (UnimplementedServiceDiscoveryServer) mustEmbedUnimplementedServiceDiscoveryServer() {}

// UnsafeServiceDiscoveryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceDiscoveryServer will
// result in compilation errors.
type UnsafeServiceDiscoveryServer interface {
	mustEmbedUnimplementedServiceDiscoveryServer()
}

func RegisterServiceDiscoveryServer(s *grpc.Server, srv ServiceDiscoveryServer) {
	s.RegisterService(&_ServiceDiscovery_serviceDesc, srv)
}

func _ServiceDiscovery_Announce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceDiscoveryServer).Announce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ServiceDiscovery/Announce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceDiscoveryServer).Announce(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServiceDiscovery_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ServiceDiscovery",
	HandlerType: (*ServiceDiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Announce",
			Handler:    _ServiceDiscovery_Announce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/announce.proto",
}