// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - ragu               v1.0.0
// source: github.com/aity-cloud/monty/pkg/plugins/apis/system/system.proto

package system

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
	System_UseManagementAPI_FullMethodName   = "/system.System/UseManagementAPI"
	System_UseKeyValueStore_FullMethodName   = "/system.System/UseKeyValueStore"
	System_UseAPIExtensions_FullMethodName   = "/system.System/UseAPIExtensions"
	System_UseCachingProvider_FullMethodName = "/system.System/UseCachingProvider"
)

// SystemClient is the client API for System service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SystemClient interface {
	UseManagementAPI(ctx context.Context, in *BrokerID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UseKeyValueStore(ctx context.Context, in *BrokerID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UseAPIExtensions(ctx context.Context, in *DialAddress, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UseCachingProvider(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type systemClient struct {
	cc grpc.ClientConnInterface
}

func NewSystemClient(cc grpc.ClientConnInterface) SystemClient {
	return &systemClient{cc}
}

func (c *systemClient) UseManagementAPI(ctx context.Context, in *BrokerID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, System_UseManagementAPI_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemClient) UseKeyValueStore(ctx context.Context, in *BrokerID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, System_UseKeyValueStore_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemClient) UseAPIExtensions(ctx context.Context, in *DialAddress, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, System_UseAPIExtensions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemClient) UseCachingProvider(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, System_UseCachingProvider_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SystemServer is the server API for System service.
// All implementations must embed UnimplementedSystemServer
// for forward compatibility
type SystemServer interface {
	UseManagementAPI(context.Context, *BrokerID) (*emptypb.Empty, error)
	UseKeyValueStore(context.Context, *BrokerID) (*emptypb.Empty, error)
	UseAPIExtensions(context.Context, *DialAddress) (*emptypb.Empty, error)
	UseCachingProvider(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedSystemServer()
}

// UnimplementedSystemServer must be embedded to have forward compatible implementations.
type UnimplementedSystemServer struct {
}

func (UnimplementedSystemServer) UseManagementAPI(context.Context, *BrokerID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseManagementAPI not implemented")
}
func (UnimplementedSystemServer) UseKeyValueStore(context.Context, *BrokerID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseKeyValueStore not implemented")
}
func (UnimplementedSystemServer) UseAPIExtensions(context.Context, *DialAddress) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseAPIExtensions not implemented")
}
func (UnimplementedSystemServer) UseCachingProvider(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseCachingProvider not implemented")
}
func (UnimplementedSystemServer) mustEmbedUnimplementedSystemServer() {}

// UnsafeSystemServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SystemServer will
// result in compilation errors.
type UnsafeSystemServer interface {
	mustEmbedUnimplementedSystemServer()
}

func RegisterSystemServer(s grpc.ServiceRegistrar, srv SystemServer) {
	s.RegisterService(&System_ServiceDesc, srv)
}

func _System_UseManagementAPI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrokerID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemServer).UseManagementAPI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: System_UseManagementAPI_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemServer).UseManagementAPI(ctx, req.(*BrokerID))
	}
	return interceptor(ctx, in, info, handler)
}

func _System_UseKeyValueStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrokerID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemServer).UseKeyValueStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: System_UseKeyValueStore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemServer).UseKeyValueStore(ctx, req.(*BrokerID))
	}
	return interceptor(ctx, in, info, handler)
}

func _System_UseAPIExtensions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DialAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemServer).UseAPIExtensions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: System_UseAPIExtensions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemServer).UseAPIExtensions(ctx, req.(*DialAddress))
	}
	return interceptor(ctx, in, info, handler)
}

func _System_UseCachingProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemServer).UseCachingProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: System_UseCachingProvider_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemServer).UseCachingProvider(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// System_ServiceDesc is the grpc.ServiceDesc for System service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var System_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "system.System",
	HandlerType: (*SystemServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UseManagementAPI",
			Handler:    _System_UseManagementAPI_Handler,
		},
		{
			MethodName: "UseKeyValueStore",
			Handler:    _System_UseKeyValueStore_Handler,
		},
		{
			MethodName: "UseAPIExtensions",
			Handler:    _System_UseAPIExtensions_Handler,
		},
		{
			MethodName: "UseCachingProvider",
			Handler:    _System_UseCachingProvider_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/pkg/plugins/apis/system/system.proto",
}

const (
	KeyValueStore_Put_FullMethodName      = "/system.KeyValueStore/Put"
	KeyValueStore_Get_FullMethodName      = "/system.KeyValueStore/Get"
	KeyValueStore_Watch_FullMethodName    = "/system.KeyValueStore/Watch"
	KeyValueStore_Delete_FullMethodName   = "/system.KeyValueStore/Delete"
	KeyValueStore_ListKeys_FullMethodName = "/system.KeyValueStore/ListKeys"
	KeyValueStore_History_FullMethodName  = "/system.KeyValueStore/History"
)

// KeyValueStoreClient is the client API for KeyValueStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeyValueStoreClient interface {
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (KeyValueStore_WatchClient, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	ListKeys(ctx context.Context, in *ListKeysRequest, opts ...grpc.CallOption) (*ListKeysResponse, error)
	History(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (*HistoryResponse, error)
}

type keyValueStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewKeyValueStoreClient(cc grpc.ClientConnInterface) KeyValueStoreClient {
	return &keyValueStoreClient{cc}
}

func (c *keyValueStoreClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, KeyValueStore_Put_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, KeyValueStore_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (KeyValueStore_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &KeyValueStore_ServiceDesc.Streams[0], KeyValueStore_Watch_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &keyValueStoreWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KeyValueStore_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type keyValueStoreWatchClient struct {
	grpc.ClientStream
}

func (x *keyValueStoreWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *keyValueStoreClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, KeyValueStore_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) ListKeys(ctx context.Context, in *ListKeysRequest, opts ...grpc.CallOption) (*ListKeysResponse, error) {
	out := new(ListKeysResponse)
	err := c.cc.Invoke(ctx, KeyValueStore_ListKeys_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) History(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (*HistoryResponse, error) {
	out := new(HistoryResponse)
	err := c.cc.Invoke(ctx, KeyValueStore_History_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyValueStoreServer is the server API for KeyValueStore service.
// All implementations must embed UnimplementedKeyValueStoreServer
// for forward compatibility
type KeyValueStoreServer interface {
	Put(context.Context, *PutRequest) (*PutResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Watch(*WatchRequest, KeyValueStore_WatchServer) error
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	ListKeys(context.Context, *ListKeysRequest) (*ListKeysResponse, error)
	History(context.Context, *HistoryRequest) (*HistoryResponse, error)
	mustEmbedUnimplementedKeyValueStoreServer()
}

// UnimplementedKeyValueStoreServer must be embedded to have forward compatible implementations.
type UnimplementedKeyValueStoreServer struct {
}

func (UnimplementedKeyValueStoreServer) Put(context.Context, *PutRequest) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedKeyValueStoreServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKeyValueStoreServer) Watch(*WatchRequest, KeyValueStore_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedKeyValueStoreServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedKeyValueStoreServer) ListKeys(context.Context, *ListKeysRequest) (*ListKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListKeys not implemented")
}
func (UnimplementedKeyValueStoreServer) History(context.Context, *HistoryRequest) (*HistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method History not implemented")
}
func (UnimplementedKeyValueStoreServer) mustEmbedUnimplementedKeyValueStoreServer() {}

// UnsafeKeyValueStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeyValueStoreServer will
// result in compilation errors.
type UnsafeKeyValueStoreServer interface {
	mustEmbedUnimplementedKeyValueStoreServer()
}

func RegisterKeyValueStoreServer(s grpc.ServiceRegistrar, srv KeyValueStoreServer) {
	s.RegisterService(&KeyValueStore_ServiceDesc, srv)
}

func _KeyValueStore_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyValueStore_Put_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyValueStore_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KeyValueStoreServer).Watch(m, &keyValueStoreWatchServer{stream})
}

type KeyValueStore_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type keyValueStoreWatchServer struct {
	grpc.ServerStream
}

func (x *keyValueStoreWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _KeyValueStore_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyValueStore_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_ListKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).ListKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyValueStore_ListKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).ListKeys(ctx, req.(*ListKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_History_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).History(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeyValueStore_History_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).History(ctx, req.(*HistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeyValueStore_ServiceDesc is the grpc.ServiceDesc for KeyValueStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeyValueStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "system.KeyValueStore",
	HandlerType: (*KeyValueStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _KeyValueStore_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KeyValueStore_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KeyValueStore_Delete_Handler,
		},
		{
			MethodName: "ListKeys",
			Handler:    _KeyValueStore_ListKeys_Handler,
		},
		{
			MethodName: "History",
			Handler:    _KeyValueStore_History_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _KeyValueStore_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/aity-cloud/monty/pkg/plugins/apis/system/system.proto",
}
