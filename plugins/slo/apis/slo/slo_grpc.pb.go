// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - ragu               v1.0.0
// source: github.com/aity-cloud/monty/plugins/slo/apis/slo/slo.proto

package slo

import (
	context "context"
	v1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
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
	SLO_CreateSLO_FullMethodName       = "/slo.SLO/CreateSLO"
	SLO_GetSLO_FullMethodName          = "/slo.SLO/GetSLO"
	SLO_ListSLOs_FullMethodName        = "/slo.SLO/ListSLOs"
	SLO_UpdateSLO_FullMethodName       = "/slo.SLO/UpdateSLO"
	SLO_DeleteSLO_FullMethodName       = "/slo.SLO/DeleteSLO"
	SLO_CloneSLO_FullMethodName        = "/slo.SLO/CloneSLO"
	SLO_CloneToClusters_FullMethodName = "/slo.SLO/CloneToClusters"
	SLO_ListMetrics_FullMethodName     = "/slo.SLO/ListMetrics"
	SLO_ListServices_FullMethodName    = "/slo.SLO/ListServices"
	SLO_ListEvents_FullMethodName      = "/slo.SLO/ListEvents"
	SLO_Status_FullMethodName          = "/slo.SLO/Status"
	SLO_Preview_FullMethodName         = "/slo.SLO/Preview"
)

// SLOClient is the client API for SLO service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SLOClient interface {
	CreateSLO(ctx context.Context, in *CreateSLORequest, opts ...grpc.CallOption) (*v1.Reference, error)
	GetSLO(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*SLOData, error)
	ListSLOs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServiceLevelObjectiveList, error)
	UpdateSLO(ctx context.Context, in *SLOData, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteSLO(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CloneSLO(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*SLOData, error)
	CloneToClusters(ctx context.Context, in *MultiClusterSLO, opts ...grpc.CallOption) (*MultiClusterFailures, error)
	// Returns a set of metrics with compatible implementations for
	// a set of services
	ListMetrics(ctx context.Context, in *ListMetricsRequest, opts ...grpc.CallOption) (*MetricGroupList, error)
	// Returns the list of services discovered by the Service Discovery backend
	ListServices(ctx context.Context, in *ListServicesRequest, opts ...grpc.CallOption) (*ServiceList, error)
	ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*EventList, error)
	// Returns a status enum badge for a given SLO
	Status(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*SLOStatus, error)
	Preview(ctx context.Context, in *CreateSLORequest, opts ...grpc.CallOption) (*SLOPreviewResponse, error)
}

type sLOClient struct {
	cc grpc.ClientConnInterface
}

func NewSLOClient(cc grpc.ClientConnInterface) SLOClient {
	return &sLOClient{cc}
}

func (c *sLOClient) CreateSLO(ctx context.Context, in *CreateSLORequest, opts ...grpc.CallOption) (*v1.Reference, error) {
	out := new(v1.Reference)
	err := c.cc.Invoke(ctx, SLO_CreateSLO_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) GetSLO(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*SLOData, error) {
	out := new(SLOData)
	err := c.cc.Invoke(ctx, SLO_GetSLO_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) ListSLOs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServiceLevelObjectiveList, error) {
	out := new(ServiceLevelObjectiveList)
	err := c.cc.Invoke(ctx, SLO_ListSLOs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) UpdateSLO(ctx context.Context, in *SLOData, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, SLO_UpdateSLO_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) DeleteSLO(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, SLO_DeleteSLO_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) CloneSLO(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*SLOData, error) {
	out := new(SLOData)
	err := c.cc.Invoke(ctx, SLO_CloneSLO_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) CloneToClusters(ctx context.Context, in *MultiClusterSLO, opts ...grpc.CallOption) (*MultiClusterFailures, error) {
	out := new(MultiClusterFailures)
	err := c.cc.Invoke(ctx, SLO_CloneToClusters_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) ListMetrics(ctx context.Context, in *ListMetricsRequest, opts ...grpc.CallOption) (*MetricGroupList, error) {
	out := new(MetricGroupList)
	err := c.cc.Invoke(ctx, SLO_ListMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) ListServices(ctx context.Context, in *ListServicesRequest, opts ...grpc.CallOption) (*ServiceList, error) {
	out := new(ServiceList)
	err := c.cc.Invoke(ctx, SLO_ListServices_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*EventList, error) {
	out := new(EventList)
	err := c.cc.Invoke(ctx, SLO_ListEvents_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) Status(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*SLOStatus, error) {
	out := new(SLOStatus)
	err := c.cc.Invoke(ctx, SLO_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLOClient) Preview(ctx context.Context, in *CreateSLORequest, opts ...grpc.CallOption) (*SLOPreviewResponse, error) {
	out := new(SLOPreviewResponse)
	err := c.cc.Invoke(ctx, SLO_Preview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SLOServer is the server API for SLO service.
// All implementations must embed UnimplementedSLOServer
// for forward compatibility
type SLOServer interface {
	CreateSLO(context.Context, *CreateSLORequest) (*v1.Reference, error)
	GetSLO(context.Context, *v1.Reference) (*SLOData, error)
	ListSLOs(context.Context, *emptypb.Empty) (*ServiceLevelObjectiveList, error)
	UpdateSLO(context.Context, *SLOData) (*emptypb.Empty, error)
	DeleteSLO(context.Context, *v1.Reference) (*emptypb.Empty, error)
	CloneSLO(context.Context, *v1.Reference) (*SLOData, error)
	CloneToClusters(context.Context, *MultiClusterSLO) (*MultiClusterFailures, error)
	// Returns a set of metrics with compatible implementations for
	// a set of services
	ListMetrics(context.Context, *ListMetricsRequest) (*MetricGroupList, error)
	// Returns the list of services discovered by the Service Discovery backend
	ListServices(context.Context, *ListServicesRequest) (*ServiceList, error)
	ListEvents(context.Context, *ListEventsRequest) (*EventList, error)
	// Returns a status enum badge for a given SLO
	Status(context.Context, *v1.Reference) (*SLOStatus, error)
	Preview(context.Context, *CreateSLORequest) (*SLOPreviewResponse, error)
	mustEmbedUnimplementedSLOServer()
}

// UnimplementedSLOServer must be embedded to have forward compatible implementations.
type UnimplementedSLOServer struct {
}

func (UnimplementedSLOServer) CreateSLO(context.Context, *CreateSLORequest) (*v1.Reference, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSLO not implemented")
}
func (UnimplementedSLOServer) GetSLO(context.Context, *v1.Reference) (*SLOData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSLO not implemented")
}
func (UnimplementedSLOServer) ListSLOs(context.Context, *emptypb.Empty) (*ServiceLevelObjectiveList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSLOs not implemented")
}
func (UnimplementedSLOServer) UpdateSLO(context.Context, *SLOData) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSLO not implemented")
}
func (UnimplementedSLOServer) DeleteSLO(context.Context, *v1.Reference) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSLO not implemented")
}
func (UnimplementedSLOServer) CloneSLO(context.Context, *v1.Reference) (*SLOData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloneSLO not implemented")
}
func (UnimplementedSLOServer) CloneToClusters(context.Context, *MultiClusterSLO) (*MultiClusterFailures, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloneToClusters not implemented")
}
func (UnimplementedSLOServer) ListMetrics(context.Context, *ListMetricsRequest) (*MetricGroupList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMetrics not implemented")
}
func (UnimplementedSLOServer) ListServices(context.Context, *ListServicesRequest) (*ServiceList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListServices not implemented")
}
func (UnimplementedSLOServer) ListEvents(context.Context, *ListEventsRequest) (*EventList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvents not implemented")
}
func (UnimplementedSLOServer) Status(context.Context, *v1.Reference) (*SLOStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedSLOServer) Preview(context.Context, *CreateSLORequest) (*SLOPreviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Preview not implemented")
}
func (UnimplementedSLOServer) mustEmbedUnimplementedSLOServer() {}

// UnsafeSLOServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SLOServer will
// result in compilation errors.
type UnsafeSLOServer interface {
	mustEmbedUnimplementedSLOServer()
}

func RegisterSLOServer(s grpc.ServiceRegistrar, srv SLOServer) {
	s.RegisterService(&SLO_ServiceDesc, srv)
}

func _SLO_CreateSLO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSLORequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).CreateSLO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_CreateSLO_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).CreateSLO(ctx, req.(*CreateSLORequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_GetSLO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Reference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).GetSLO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_GetSLO_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).GetSLO(ctx, req.(*v1.Reference))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_ListSLOs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).ListSLOs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_ListSLOs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).ListSLOs(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_UpdateSLO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SLOData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).UpdateSLO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_UpdateSLO_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).UpdateSLO(ctx, req.(*SLOData))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_DeleteSLO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Reference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).DeleteSLO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_DeleteSLO_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).DeleteSLO(ctx, req.(*v1.Reference))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_CloneSLO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Reference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).CloneSLO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_CloneSLO_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).CloneSLO(ctx, req.(*v1.Reference))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_CloneToClusters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiClusterSLO)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).CloneToClusters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_CloneToClusters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).CloneToClusters(ctx, req.(*MultiClusterSLO))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_ListMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).ListMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_ListMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).ListMetrics(ctx, req.(*ListMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_ListServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).ListServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_ListServices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).ListServices(ctx, req.(*ListServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).ListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_ListEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).ListEvents(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Reference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).Status(ctx, req.(*v1.Reference))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLO_Preview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSLORequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLOServer).Preview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLO_Preview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLOServer).Preview(ctx, req.(*CreateSLORequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SLO_ServiceDesc is the grpc.ServiceDesc for SLO service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SLO_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "slo.SLO",
	HandlerType: (*SLOServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSLO",
			Handler:    _SLO_CreateSLO_Handler,
		},
		{
			MethodName: "GetSLO",
			Handler:    _SLO_GetSLO_Handler,
		},
		{
			MethodName: "ListSLOs",
			Handler:    _SLO_ListSLOs_Handler,
		},
		{
			MethodName: "UpdateSLO",
			Handler:    _SLO_UpdateSLO_Handler,
		},
		{
			MethodName: "DeleteSLO",
			Handler:    _SLO_DeleteSLO_Handler,
		},
		{
			MethodName: "CloneSLO",
			Handler:    _SLO_CloneSLO_Handler,
		},
		{
			MethodName: "CloneToClusters",
			Handler:    _SLO_CloneToClusters_Handler,
		},
		{
			MethodName: "ListMetrics",
			Handler:    _SLO_ListMetrics_Handler,
		},
		{
			MethodName: "ListServices",
			Handler:    _SLO_ListServices_Handler,
		},
		{
			MethodName: "ListEvents",
			Handler:    _SLO_ListEvents_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _SLO_Status_Handler,
		},
		{
			MethodName: "Preview",
			Handler:    _SLO_Preview_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/plugins/slo/apis/slo/slo.proto",
}
