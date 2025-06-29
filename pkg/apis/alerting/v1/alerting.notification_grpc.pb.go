// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - ragu               v1.0.0
// source: github.com/aity-cloud/monty/pkg/apis/alerting/v1/alerting.notification.proto

package v1

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
	AlertNotifications_TestAlertEndpoint_FullMethodName        = "/alerting.AlertNotifications/TestAlertEndpoint"
	AlertNotifications_TriggerAlerts_FullMethodName            = "/alerting.AlertNotifications/TriggerAlerts"
	AlertNotifications_ResolveAlerts_FullMethodName            = "/alerting.AlertNotifications/ResolveAlerts"
	AlertNotifications_PushNotification_FullMethodName         = "/alerting.AlertNotifications/PushNotification"
	AlertNotifications_ListNotifications_FullMethodName        = "/alerting.AlertNotifications/ListNotifications"
	AlertNotifications_ListAlarmMessages_FullMethodName        = "/alerting.AlertNotifications/ListAlarmMessages"
	AlertNotifications_ListRoutingRelationships_FullMethodName = "/alerting.AlertNotifications/ListRoutingRelationships"
)

// AlertNotificationsClient is the client API for AlertNotifications service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AlertNotificationsClient interface {
	TestAlertEndpoint(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*emptypb.Empty, error)
	TriggerAlerts(ctx context.Context, in *TriggerAlertsRequest, opts ...grpc.CallOption) (*TriggerAlertsResponse, error)
	ResolveAlerts(ctx context.Context, in *ResolveAlertsRequest, opts ...grpc.CallOption) (*ResolveAlertsResponse, error)
	PushNotification(ctx context.Context, in *Notification, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// In the cache we evict the keys with the highest (priority,severity)
	// according to the given filter
	// but return the filtered messages sorted by timestamp.
	ListNotifications(ctx context.Context, in *ListNotificationRequest, opts ...grpc.CallOption) (*ListMessageResponse, error)
	// best-effort listing of alarm messages for a given window
	// messages with low frequency and severity are dropped frequently
	// so may not show up with their associated incident
	ListAlarmMessages(ctx context.Context, in *ListAlarmMessageRequest, opts ...grpc.CallOption) (*ListMessageResponse, error)
	ListRoutingRelationships(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListRoutingRelationshipsResponse, error)
}

type alertNotificationsClient struct {
	cc grpc.ClientConnInterface
}

func NewAlertNotificationsClient(cc grpc.ClientConnInterface) AlertNotificationsClient {
	return &alertNotificationsClient{cc}
}

func (c *alertNotificationsClient) TestAlertEndpoint(ctx context.Context, in *v1.Reference, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AlertNotifications_TestAlertEndpoint_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertNotificationsClient) TriggerAlerts(ctx context.Context, in *TriggerAlertsRequest, opts ...grpc.CallOption) (*TriggerAlertsResponse, error) {
	out := new(TriggerAlertsResponse)
	err := c.cc.Invoke(ctx, AlertNotifications_TriggerAlerts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertNotificationsClient) ResolveAlerts(ctx context.Context, in *ResolveAlertsRequest, opts ...grpc.CallOption) (*ResolveAlertsResponse, error) {
	out := new(ResolveAlertsResponse)
	err := c.cc.Invoke(ctx, AlertNotifications_ResolveAlerts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertNotificationsClient) PushNotification(ctx context.Context, in *Notification, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AlertNotifications_PushNotification_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertNotificationsClient) ListNotifications(ctx context.Context, in *ListNotificationRequest, opts ...grpc.CallOption) (*ListMessageResponse, error) {
	out := new(ListMessageResponse)
	err := c.cc.Invoke(ctx, AlertNotifications_ListNotifications_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertNotificationsClient) ListAlarmMessages(ctx context.Context, in *ListAlarmMessageRequest, opts ...grpc.CallOption) (*ListMessageResponse, error) {
	out := new(ListMessageResponse)
	err := c.cc.Invoke(ctx, AlertNotifications_ListAlarmMessages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertNotificationsClient) ListRoutingRelationships(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListRoutingRelationshipsResponse, error) {
	out := new(ListRoutingRelationshipsResponse)
	err := c.cc.Invoke(ctx, AlertNotifications_ListRoutingRelationships_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlertNotificationsServer is the server API for AlertNotifications service.
// All implementations must embed UnimplementedAlertNotificationsServer
// for forward compatibility
type AlertNotificationsServer interface {
	TestAlertEndpoint(context.Context, *v1.Reference) (*emptypb.Empty, error)
	TriggerAlerts(context.Context, *TriggerAlertsRequest) (*TriggerAlertsResponse, error)
	ResolveAlerts(context.Context, *ResolveAlertsRequest) (*ResolveAlertsResponse, error)
	PushNotification(context.Context, *Notification) (*emptypb.Empty, error)
	// In the cache we evict the keys with the highest (priority,severity)
	// according to the given filter
	// but return the filtered messages sorted by timestamp.
	ListNotifications(context.Context, *ListNotificationRequest) (*ListMessageResponse, error)
	// best-effort listing of alarm messages for a given window
	// messages with low frequency and severity are dropped frequently
	// so may not show up with their associated incident
	ListAlarmMessages(context.Context, *ListAlarmMessageRequest) (*ListMessageResponse, error)
	ListRoutingRelationships(context.Context, *emptypb.Empty) (*ListRoutingRelationshipsResponse, error)
	mustEmbedUnimplementedAlertNotificationsServer()
}

// UnimplementedAlertNotificationsServer must be embedded to have forward compatible implementations.
type UnimplementedAlertNotificationsServer struct {
}

func (UnimplementedAlertNotificationsServer) TestAlertEndpoint(context.Context, *v1.Reference) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestAlertEndpoint not implemented")
}
func (UnimplementedAlertNotificationsServer) TriggerAlerts(context.Context, *TriggerAlertsRequest) (*TriggerAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerAlerts not implemented")
}
func (UnimplementedAlertNotificationsServer) ResolveAlerts(context.Context, *ResolveAlertsRequest) (*ResolveAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveAlerts not implemented")
}
func (UnimplementedAlertNotificationsServer) PushNotification(context.Context, *Notification) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushNotification not implemented")
}
func (UnimplementedAlertNotificationsServer) ListNotifications(context.Context, *ListNotificationRequest) (*ListMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNotifications not implemented")
}
func (UnimplementedAlertNotificationsServer) ListAlarmMessages(context.Context, *ListAlarmMessageRequest) (*ListMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAlarmMessages not implemented")
}
func (UnimplementedAlertNotificationsServer) ListRoutingRelationships(context.Context, *emptypb.Empty) (*ListRoutingRelationshipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoutingRelationships not implemented")
}
func (UnimplementedAlertNotificationsServer) mustEmbedUnimplementedAlertNotificationsServer() {}

// UnsafeAlertNotificationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AlertNotificationsServer will
// result in compilation errors.
type UnsafeAlertNotificationsServer interface {
	mustEmbedUnimplementedAlertNotificationsServer()
}

func RegisterAlertNotificationsServer(s grpc.ServiceRegistrar, srv AlertNotificationsServer) {
	s.RegisterService(&AlertNotifications_ServiceDesc, srv)
}

func _AlertNotifications_TestAlertEndpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Reference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertNotificationsServer).TestAlertEndpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertNotifications_TestAlertEndpoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertNotificationsServer).TestAlertEndpoint(ctx, req.(*v1.Reference))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertNotifications_TriggerAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertNotificationsServer).TriggerAlerts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertNotifications_TriggerAlerts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertNotificationsServer).TriggerAlerts(ctx, req.(*TriggerAlertsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertNotifications_ResolveAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertNotificationsServer).ResolveAlerts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertNotifications_ResolveAlerts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertNotificationsServer).ResolveAlerts(ctx, req.(*ResolveAlertsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertNotifications_PushNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Notification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertNotificationsServer).PushNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertNotifications_PushNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertNotificationsServer).PushNotification(ctx, req.(*Notification))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertNotifications_ListNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertNotificationsServer).ListNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertNotifications_ListNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertNotificationsServer).ListNotifications(ctx, req.(*ListNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertNotifications_ListAlarmMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAlarmMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertNotificationsServer).ListAlarmMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertNotifications_ListAlarmMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertNotificationsServer).ListAlarmMessages(ctx, req.(*ListAlarmMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertNotifications_ListRoutingRelationships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertNotificationsServer).ListRoutingRelationships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertNotifications_ListRoutingRelationships_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertNotificationsServer).ListRoutingRelationships(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// AlertNotifications_ServiceDesc is the grpc.ServiceDesc for AlertNotifications service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AlertNotifications_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "alerting.AlertNotifications",
	HandlerType: (*AlertNotificationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TestAlertEndpoint",
			Handler:    _AlertNotifications_TestAlertEndpoint_Handler,
		},
		{
			MethodName: "TriggerAlerts",
			Handler:    _AlertNotifications_TriggerAlerts_Handler,
		},
		{
			MethodName: "ResolveAlerts",
			Handler:    _AlertNotifications_ResolveAlerts_Handler,
		},
		{
			MethodName: "PushNotification",
			Handler:    _AlertNotifications_PushNotification_Handler,
		},
		{
			MethodName: "ListNotifications",
			Handler:    _AlertNotifications_ListNotifications_Handler,
		},
		{
			MethodName: "ListAlarmMessages",
			Handler:    _AlertNotifications_ListAlarmMessages_Handler,
		},
		{
			MethodName: "ListRoutingRelationships",
			Handler:    _AlertNotifications_ListRoutingRelationships_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/pkg/apis/alerting/v1/alerting.notification.proto",
}
