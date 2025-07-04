// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - ragu               v1.0.0
// source: github.com/aity-cloud/monty/plugins/example/pkg/example/example.proto

package example

import (
	context "context"
	driverutil "github.com/aity-cloud/monty/pkg/plugins/driverutil"
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
	ExampleAPIExtension_Echo_FullMethodName  = "/example.ExampleAPIExtension/Echo"
	ExampleAPIExtension_Ready_FullMethodName = "/example.ExampleAPIExtension/Ready"
)

// ExampleAPIExtensionClient is the client API for ExampleAPIExtension service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExampleAPIExtensionClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	Ready(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type exampleAPIExtensionClient struct {
	cc grpc.ClientConnInterface
}

func NewExampleAPIExtensionClient(cc grpc.ClientConnInterface) ExampleAPIExtensionClient {
	return &exampleAPIExtensionClient{cc}
}

func (c *exampleAPIExtensionClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, ExampleAPIExtension_Echo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleAPIExtensionClient) Ready(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ExampleAPIExtension_Ready_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExampleAPIExtensionServer is the server API for ExampleAPIExtension service.
// All implementations must embed UnimplementedExampleAPIExtensionServer
// for forward compatibility
type ExampleAPIExtensionServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	Ready(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedExampleAPIExtensionServer()
}

// UnimplementedExampleAPIExtensionServer must be embedded to have forward compatible implementations.
type UnimplementedExampleAPIExtensionServer struct {
}

func (UnimplementedExampleAPIExtensionServer) Echo(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedExampleAPIExtensionServer) Ready(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ready not implemented")
}
func (UnimplementedExampleAPIExtensionServer) mustEmbedUnimplementedExampleAPIExtensionServer() {}

// UnsafeExampleAPIExtensionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExampleAPIExtensionServer will
// result in compilation errors.
type UnsafeExampleAPIExtensionServer interface {
	mustEmbedUnimplementedExampleAPIExtensionServer()
}

func RegisterExampleAPIExtensionServer(s grpc.ServiceRegistrar, srv ExampleAPIExtensionServer) {
	s.RegisterService(&ExampleAPIExtension_ServiceDesc, srv)
}

func _ExampleAPIExtension_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleAPIExtensionServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExampleAPIExtension_Echo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleAPIExtensionServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExampleAPIExtension_Ready_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleAPIExtensionServer).Ready(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExampleAPIExtension_Ready_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleAPIExtensionServer).Ready(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ExampleAPIExtension_ServiceDesc is the grpc.ServiceDesc for ExampleAPIExtension service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExampleAPIExtension_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "example.ExampleAPIExtension",
	HandlerType: (*ExampleAPIExtensionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _ExampleAPIExtension_Echo_Handler,
		},
		{
			MethodName: "Ready",
			Handler:    _ExampleAPIExtension_Ready_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/plugins/example/pkg/example/example.proto",
}

const (
	ExampleUnaryExtension_Hello_FullMethodName = "/example.ExampleUnaryExtension/Hello"
)

// ExampleUnaryExtensionClient is the client API for ExampleUnaryExtension service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExampleUnaryExtensionClient interface {
	Hello(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*EchoResponse, error)
}

type exampleUnaryExtensionClient struct {
	cc grpc.ClientConnInterface
}

func NewExampleUnaryExtensionClient(cc grpc.ClientConnInterface) ExampleUnaryExtensionClient {
	return &exampleUnaryExtensionClient{cc}
}

func (c *exampleUnaryExtensionClient) Hello(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, ExampleUnaryExtension_Hello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExampleUnaryExtensionServer is the server API for ExampleUnaryExtension service.
// All implementations must embed UnimplementedExampleUnaryExtensionServer
// for forward compatibility
type ExampleUnaryExtensionServer interface {
	Hello(context.Context, *emptypb.Empty) (*EchoResponse, error)
	mustEmbedUnimplementedExampleUnaryExtensionServer()
}

// UnimplementedExampleUnaryExtensionServer must be embedded to have forward compatible implementations.
type UnimplementedExampleUnaryExtensionServer struct {
}

func (UnimplementedExampleUnaryExtensionServer) Hello(context.Context, *emptypb.Empty) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedExampleUnaryExtensionServer) mustEmbedUnimplementedExampleUnaryExtensionServer() {}

// UnsafeExampleUnaryExtensionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExampleUnaryExtensionServer will
// result in compilation errors.
type UnsafeExampleUnaryExtensionServer interface {
	mustEmbedUnimplementedExampleUnaryExtensionServer()
}

func RegisterExampleUnaryExtensionServer(s grpc.ServiceRegistrar, srv ExampleUnaryExtensionServer) {
	s.RegisterService(&ExampleUnaryExtension_ServiceDesc, srv)
}

func _ExampleUnaryExtension_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleUnaryExtensionServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExampleUnaryExtension_Hello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleUnaryExtensionServer).Hello(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ExampleUnaryExtension_ServiceDesc is the grpc.ServiceDesc for ExampleUnaryExtension service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExampleUnaryExtension_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "example.ExampleUnaryExtension",
	HandlerType: (*ExampleUnaryExtensionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _ExampleUnaryExtension_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/plugins/example/pkg/example/example.proto",
}

const (
	Config_GetDefaultConfiguration_FullMethodName   = "/example.Config/GetDefaultConfiguration"
	Config_SetDefaultConfiguration_FullMethodName   = "/example.Config/SetDefaultConfiguration"
	Config_GetConfiguration_FullMethodName          = "/example.Config/GetConfiguration"
	Config_SetConfiguration_FullMethodName          = "/example.Config/SetConfiguration"
	Config_ResetDefaultConfiguration_FullMethodName = "/example.Config/ResetDefaultConfiguration"
	Config_ResetConfiguration_FullMethodName        = "/example.Config/ResetConfiguration"
	Config_DryRun_FullMethodName                    = "/example.Config/DryRun"
	Config_ConfigurationHistory_FullMethodName      = "/example.Config/ConfigurationHistory"
)

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigClient interface {
	GetDefaultConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*ConfigSpec, error)
	SetDefaultConfiguration(ctx context.Context, in *ConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*ConfigSpec, error)
	SetConfiguration(ctx context.Context, in *ConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ResetDefaultConfiguration(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ResetConfiguration(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DryRun(ctx context.Context, in *DryRunRequest, opts ...grpc.CallOption) (*DryRunResponse, error)
	ConfigurationHistory(ctx context.Context, in *driverutil.ConfigurationHistoryRequest, opts ...grpc.CallOption) (*HistoryResponse, error)
}

type configClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigClient(cc grpc.ClientConnInterface) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) GetDefaultConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*ConfigSpec, error) {
	out := new(ConfigSpec)
	err := c.cc.Invoke(ctx, Config_GetDefaultConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) SetDefaultConfiguration(ctx context.Context, in *ConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Config_SetDefaultConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) GetConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*ConfigSpec, error) {
	out := new(ConfigSpec)
	err := c.cc.Invoke(ctx, Config_GetConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) SetConfiguration(ctx context.Context, in *ConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Config_SetConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) ResetDefaultConfiguration(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Config_ResetDefaultConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) ResetConfiguration(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Config_ResetConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) DryRun(ctx context.Context, in *DryRunRequest, opts ...grpc.CallOption) (*DryRunResponse, error) {
	out := new(DryRunResponse)
	err := c.cc.Invoke(ctx, Config_DryRun_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) ConfigurationHistory(ctx context.Context, in *driverutil.ConfigurationHistoryRequest, opts ...grpc.CallOption) (*HistoryResponse, error) {
	out := new(HistoryResponse)
	err := c.cc.Invoke(ctx, Config_ConfigurationHistory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
// All implementations must embed UnimplementedConfigServer
// for forward compatibility
type ConfigServer interface {
	GetDefaultConfiguration(context.Context, *driverutil.GetRequest) (*ConfigSpec, error)
	SetDefaultConfiguration(context.Context, *ConfigSpec) (*emptypb.Empty, error)
	GetConfiguration(context.Context, *driverutil.GetRequest) (*ConfigSpec, error)
	SetConfiguration(context.Context, *ConfigSpec) (*emptypb.Empty, error)
	ResetDefaultConfiguration(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	ResetConfiguration(context.Context, *ResetRequest) (*emptypb.Empty, error)
	DryRun(context.Context, *DryRunRequest) (*DryRunResponse, error)
	ConfigurationHistory(context.Context, *driverutil.ConfigurationHistoryRequest) (*HistoryResponse, error)
	mustEmbedUnimplementedConfigServer()
}

// UnimplementedConfigServer must be embedded to have forward compatible implementations.
type UnimplementedConfigServer struct {
}

func (UnimplementedConfigServer) GetDefaultConfiguration(context.Context, *driverutil.GetRequest) (*ConfigSpec, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultConfiguration not implemented")
}
func (UnimplementedConfigServer) SetDefaultConfiguration(context.Context, *ConfigSpec) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDefaultConfiguration not implemented")
}
func (UnimplementedConfigServer) GetConfiguration(context.Context, *driverutil.GetRequest) (*ConfigSpec, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfiguration not implemented")
}
func (UnimplementedConfigServer) SetConfiguration(context.Context, *ConfigSpec) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetConfiguration not implemented")
}
func (UnimplementedConfigServer) ResetDefaultConfiguration(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetDefaultConfiguration not implemented")
}
func (UnimplementedConfigServer) ResetConfiguration(context.Context, *ResetRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetConfiguration not implemented")
}
func (UnimplementedConfigServer) DryRun(context.Context, *DryRunRequest) (*DryRunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DryRun not implemented")
}
func (UnimplementedConfigServer) ConfigurationHistory(context.Context, *driverutil.ConfigurationHistoryRequest) (*HistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigurationHistory not implemented")
}
func (UnimplementedConfigServer) mustEmbedUnimplementedConfigServer() {}

// UnsafeConfigServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigServer will
// result in compilation errors.
type UnsafeConfigServer interface {
	mustEmbedUnimplementedConfigServer()
}

func RegisterConfigServer(s grpc.ServiceRegistrar, srv ConfigServer) {
	s.RegisterService(&Config_ServiceDesc, srv)
}

func _Config_GetDefaultConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(driverutil.GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetDefaultConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_GetDefaultConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetDefaultConfiguration(ctx, req.(*driverutil.GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_SetDefaultConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).SetDefaultConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_SetDefaultConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).SetDefaultConfiguration(ctx, req.(*ConfigSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_GetConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(driverutil.GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_GetConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetConfiguration(ctx, req.(*driverutil.GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_SetConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).SetConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_SetConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).SetConfiguration(ctx, req.(*ConfigSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_ResetDefaultConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).ResetDefaultConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_ResetDefaultConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).ResetDefaultConfiguration(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_ResetConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).ResetConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_ResetConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).ResetConfiguration(ctx, req.(*ResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_DryRun_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DryRunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).DryRun(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_DryRun_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).DryRun(ctx, req.(*DryRunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_ConfigurationHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(driverutil.ConfigurationHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).ConfigurationHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_ConfigurationHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).ConfigurationHistory(ctx, req.(*driverutil.ConfigurationHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Config_ServiceDesc is the grpc.ServiceDesc for Config service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Config_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "example.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDefaultConfiguration",
			Handler:    _Config_GetDefaultConfiguration_Handler,
		},
		{
			MethodName: "SetDefaultConfiguration",
			Handler:    _Config_SetDefaultConfiguration_Handler,
		},
		{
			MethodName: "GetConfiguration",
			Handler:    _Config_GetConfiguration_Handler,
		},
		{
			MethodName: "SetConfiguration",
			Handler:    _Config_SetConfiguration_Handler,
		},
		{
			MethodName: "ResetDefaultConfiguration",
			Handler:    _Config_ResetDefaultConfiguration_Handler,
		},
		{
			MethodName: "ResetConfiguration",
			Handler:    _Config_ResetConfiguration_Handler,
		},
		{
			MethodName: "DryRun",
			Handler:    _Config_DryRun_Handler,
		},
		{
			MethodName: "ConfigurationHistory",
			Handler:    _Config_ConfigurationHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/plugins/example/pkg/example/example.proto",
}
