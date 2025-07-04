// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - ragu               v1.0.0
// source: github.com/aity-cloud/monty/plugins/topology/apis/orchestrator/orchestrator.proto

package orchestrator

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

const (
	TopologyOrchestrator_GetClusterStatus_FullMethodName = "/orchestrator.TopologyOrchestrator/GetClusterStatus"
)

// TopologyOrchestratorClient is the client API for TopologyOrchestrator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TopologyOrchestratorClient interface {
	GetClusterStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*InstallStatus, error)
}

type topologyOrchestratorClient struct {
	cc grpc.ClientConnInterface
}

func NewTopologyOrchestratorClient(cc grpc.ClientConnInterface) TopologyOrchestratorClient {
	return &topologyOrchestratorClient{cc}
}

func (c *topologyOrchestratorClient) GetClusterStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*InstallStatus, error) {
	out := new(InstallStatus)
	err := c.cc.Invoke(ctx, TopologyOrchestrator_GetClusterStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TopologyOrchestratorServer is the server API for TopologyOrchestrator service.
// All implementations must embed UnimplementedTopologyOrchestratorServer
// for forward compatibility
type TopologyOrchestratorServer interface {
	GetClusterStatus(context.Context, *emptypb.Empty) (*InstallStatus, error)
	mustEmbedUnimplementedTopologyOrchestratorServer()
}

// UnimplementedTopologyOrchestratorServer must be embedded to have forward compatible implementations.
type UnimplementedTopologyOrchestratorServer struct {
}

func (UnimplementedTopologyOrchestratorServer) GetClusterStatus(context.Context, *emptypb.Empty) (*InstallStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterStatus not implemented")
}
func (UnimplementedTopologyOrchestratorServer) mustEmbedUnimplementedTopologyOrchestratorServer() {}

// UnsafeTopologyOrchestratorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TopologyOrchestratorServer will
// result in compilation errors.
type UnsafeTopologyOrchestratorServer interface {
	mustEmbedUnimplementedTopologyOrchestratorServer()
}

func RegisterTopologyOrchestratorServer(s grpc.ServiceRegistrar, srv TopologyOrchestratorServer) {
	s.RegisterService(&TopologyOrchestrator_ServiceDesc, srv)
}

func _TopologyOrchestrator_GetClusterStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopologyOrchestratorServer).GetClusterStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TopologyOrchestrator_GetClusterStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopologyOrchestratorServer).GetClusterStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// TopologyOrchestrator_ServiceDesc is the grpc.ServiceDesc for TopologyOrchestrator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TopologyOrchestrator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "orchestrator.TopologyOrchestrator",
	HandlerType: (*TopologyOrchestratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetClusterStatus",
			Handler:    _TopologyOrchestrator_GetClusterStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/plugins/topology/apis/orchestrator/orchestrator.proto",
}
