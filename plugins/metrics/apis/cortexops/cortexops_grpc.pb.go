// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - ragu               v1.0.0
// source: github.com/aity-cloud/monty/plugins/metrics/apis/cortexops/cortexops.proto

package cortexops

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
	CortexOps_GetDefaultConfiguration_FullMethodName   = "/cortexops.CortexOps/GetDefaultConfiguration"
	CortexOps_SetDefaultConfiguration_FullMethodName   = "/cortexops.CortexOps/SetDefaultConfiguration"
	CortexOps_ResetDefaultConfiguration_FullMethodName = "/cortexops.CortexOps/ResetDefaultConfiguration"
	CortexOps_GetConfiguration_FullMethodName          = "/cortexops.CortexOps/GetConfiguration"
	CortexOps_SetConfiguration_FullMethodName          = "/cortexops.CortexOps/SetConfiguration"
	CortexOps_ResetConfiguration_FullMethodName        = "/cortexops.CortexOps/ResetConfiguration"
	CortexOps_Status_FullMethodName                    = "/cortexops.CortexOps/Status"
	CortexOps_Install_FullMethodName                   = "/cortexops.CortexOps/Install"
	CortexOps_Uninstall_FullMethodName                 = "/cortexops.CortexOps/Uninstall"
	CortexOps_ListPresets_FullMethodName               = "/cortexops.CortexOps/ListPresets"
	CortexOps_DryRun_FullMethodName                    = "/cortexops.CortexOps/DryRun"
	CortexOps_ConfigurationHistory_FullMethodName      = "/cortexops.CortexOps/ConfigurationHistory"
)

// CortexOpsClient is the client API for CortexOps service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CortexOpsClient interface {
	// Returns the default implementation-specific configuration, or one previously set.
	// If a default configuration was previously set using SetDefaultConfiguration, it
	// returns that configuration. Otherwise, returns implementation-specific defaults.
	// An optional revision argument can be provided to get a specific historical
	// version of the configuration instead of the current configuration.
	GetDefaultConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*CapabilityBackendConfigSpec, error)
	// Sets the default configuration that will be used as the base for future configuration changes.
	// If no custom default configuration is set using this method,
	// implementation-specific defaults may be chosen.
	// If all fields are unset, this will clear any previously-set default configuration
	// and revert back to the implementation-specific defaults.
	//
	// This API is different from the SetConfiguration API, and should not be necessary
	// for most use cases. It can be used in situations where an additional persistence
	// layer that is not driver-specific is desired.
	SetDefaultConfiguration(ctx context.Context, in *CapabilityBackendConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Resets the default configuration to the implementation-specific defaults.
	// If a custom default configuration was previously set using SetDefaultConfiguration,
	// this will clear it and revert back to the implementation-specific defaults.
	// Otherwise, this will have no effect.
	ResetDefaultConfiguration(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Gets the current configuration of the managed Cortex cluster.
	// An optional revision argument can be provided to get a specific historical
	// version of the configuration instead of the current configuration.
	GetConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*CapabilityBackendConfigSpec, error)
	// Updates the configuration of the managed Cortex cluster to match the provided configuration.
	// If the cluster is not installed, it will be configured, but remain disabled.
	// Otherwise, the already-installed cluster will be reconfigured.
	// The provided configuration will be merged with the default configuration
	// by directly overwriting fields. Slices and maps are overwritten and not combined.
	// Subsequent calls to this API will merge inputs with the current configuration,
	// not the default configuration.
	// When updating an existing configuration, the revision number in the updated configuration
	// must match the revision number of the existing configuration, otherwise a conflict
	// error will be returned. The timestamp field of the revision is ignored.
	//
	// Note: some fields may contain secrets. The placeholder value "***" can be used to
	// keep an existing secret when updating the cluster configuration.
	SetConfiguration(ctx context.Context, in *CapabilityBackendConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Resets the configuration of the managed Cortex cluster to the current default configuration.
	//
	// The request may optionally contain a field mask to specify which fields should
	// be preserved. Furthermore, if a mask is set, the request may also contain a patch
	// object used to apply additional changes to the masked fields. These changes are
	// applied atomically at the time of reset. Fields present in the patch object, but
	// not in the mask, are ignored.
	//
	// For example, with the following message:
	//
	//	message Example {
	//	  optional int32 a = 1;
	//	  optional int32 b = 2;
	//	  optional int32 c = 3;
	//	}
	//
	// and current state:
	//
	//	active:  { a: 1, b: 2, c: 3 }
	//	default: { a: 4, b: 5, c: 6 }
	//
	// and reset request parameters:
	//
	//	{
	//	  mask:    { paths: [ "a", "b" ] }
	//	  patch:   { a: 100 }
	//	}
	//
	// The resulting active configuration will be:
	//
	//	active:  {
	//	  a: 100, // masked, set to 100 via patch
	//	  b: 2,   // masked, but not set in patch, so left unchanged
	//	  c: 6,   // not masked, reset to default
	//	}
	ResetConfiguration(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Gets the current status of the managed Cortex cluster.
	// The status includes the current install state, version, and metadata. If
	// the cluster is in the process of being reconfigured or uninstalled, it will
	// be reflected in the install state.
	// No guarantees are made about the contents of the metadata field; its
	// contents are strictly informational.
	Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*driverutil.InstallStatus, error)
	// Installs the managed Cortex cluster.
	// The cluster will be installed using the current configuration, or the
	// default configuration if none is explicitly set.
	Install(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Uninstalls the managed Cortex cluster.
	// Implementation details including error handling and system state requirements
	// are left to the cluster driver, and this API makes no guarantees about
	// the state of the cluster after the call completes (regardless of success).
	Uninstall(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Returns a static list of presets that can be used as a base for configuring the managed Cortex cluster.
	// There are several ways to use the presets, depending
	// on the desired behavior:
	//  1. Set the default configuration to a preset spec, then use SetConfiguration
	//     to fill in any additional required fields (credentials, etc)
	//  2. Add the required fields to the default configuration, then use
	//     SetConfiguration with a preset spec.
	//  3. Leave the default configuration as-is, and use SetConfiguration with a
	//     preset spec plus the required fields.
	ListPresets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PresetList, error)
	// Show what changes would be made to a configuration without saving them.
	// The request expects an action, target, and spec to be provided. These
	// correspond roughly to the other APIs in this service.
	//
	// Configuring DryRunRequest:
	//   - Use the Active target for the SetConfiguration API, and the Default target
	//     for the SetDefaultConfiguration API. Install and Uninstall actions do not
	//     require a target.
	//   - Only the Set action requires a spec to be provided.
	//
	// Notes:
	//   - When DryRun is used on Install or Uninstall requests, the response will
	//     contain modifications to the 'enabled' field only. This field is read-only
	//     in the Set* APIs.
	//   - To validate the current configuration but keep it unchanged, use the
	//     Set action with an empty spec.
	//   - Configurations returned by DryRun will always have an empty revision field.
	DryRun(ctx context.Context, in *DryRunRequest, opts ...grpc.CallOption) (*DryRunResponse, error)
	// Get a list of all past revisions of the configuration.
	// Will return the history for either the active or default configuration
	// depending on the specified target.
	// The entries are ordered from oldest to newest, where the last entry is
	// the current configuration.
	ConfigurationHistory(ctx context.Context, in *driverutil.ConfigurationHistoryRequest, opts ...grpc.CallOption) (*ConfigurationHistoryResponse, error)
}

type cortexOpsClient struct {
	cc grpc.ClientConnInterface
}

func NewCortexOpsClient(cc grpc.ClientConnInterface) CortexOpsClient {
	return &cortexOpsClient{cc}
}

func (c *cortexOpsClient) GetDefaultConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*CapabilityBackendConfigSpec, error) {
	out := new(CapabilityBackendConfigSpec)
	err := c.cc.Invoke(ctx, CortexOps_GetDefaultConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) SetDefaultConfiguration(ctx context.Context, in *CapabilityBackendConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CortexOps_SetDefaultConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) ResetDefaultConfiguration(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CortexOps_ResetDefaultConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) GetConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*CapabilityBackendConfigSpec, error) {
	out := new(CapabilityBackendConfigSpec)
	err := c.cc.Invoke(ctx, CortexOps_GetConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) SetConfiguration(ctx context.Context, in *CapabilityBackendConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CortexOps_SetConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) ResetConfiguration(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CortexOps_ResetConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*driverutil.InstallStatus, error) {
	out := new(driverutil.InstallStatus)
	err := c.cc.Invoke(ctx, CortexOps_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) Install(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CortexOps_Install_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) Uninstall(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CortexOps_Uninstall_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) ListPresets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PresetList, error) {
	out := new(PresetList)
	err := c.cc.Invoke(ctx, CortexOps_ListPresets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) DryRun(ctx context.Context, in *DryRunRequest, opts ...grpc.CallOption) (*DryRunResponse, error) {
	out := new(DryRunResponse)
	err := c.cc.Invoke(ctx, CortexOps_DryRun_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cortexOpsClient) ConfigurationHistory(ctx context.Context, in *driverutil.ConfigurationHistoryRequest, opts ...grpc.CallOption) (*ConfigurationHistoryResponse, error) {
	out := new(ConfigurationHistoryResponse)
	err := c.cc.Invoke(ctx, CortexOps_ConfigurationHistory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CortexOpsServer is the server API for CortexOps service.
// All implementations must embed UnimplementedCortexOpsServer
// for forward compatibility
type CortexOpsServer interface {
	// Returns the default implementation-specific configuration, or one previously set.
	// If a default configuration was previously set using SetDefaultConfiguration, it
	// returns that configuration. Otherwise, returns implementation-specific defaults.
	// An optional revision argument can be provided to get a specific historical
	// version of the configuration instead of the current configuration.
	GetDefaultConfiguration(context.Context, *driverutil.GetRequest) (*CapabilityBackendConfigSpec, error)
	// Sets the default configuration that will be used as the base for future configuration changes.
	// If no custom default configuration is set using this method,
	// implementation-specific defaults may be chosen.
	// If all fields are unset, this will clear any previously-set default configuration
	// and revert back to the implementation-specific defaults.
	//
	// This API is different from the SetConfiguration API, and should not be necessary
	// for most use cases. It can be used in situations where an additional persistence
	// layer that is not driver-specific is desired.
	SetDefaultConfiguration(context.Context, *CapabilityBackendConfigSpec) (*emptypb.Empty, error)
	// Resets the default configuration to the implementation-specific defaults.
	// If a custom default configuration was previously set using SetDefaultConfiguration,
	// this will clear it and revert back to the implementation-specific defaults.
	// Otherwise, this will have no effect.
	ResetDefaultConfiguration(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// Gets the current configuration of the managed Cortex cluster.
	// An optional revision argument can be provided to get a specific historical
	// version of the configuration instead of the current configuration.
	GetConfiguration(context.Context, *driverutil.GetRequest) (*CapabilityBackendConfigSpec, error)
	// Updates the configuration of the managed Cortex cluster to match the provided configuration.
	// If the cluster is not installed, it will be configured, but remain disabled.
	// Otherwise, the already-installed cluster will be reconfigured.
	// The provided configuration will be merged with the default configuration
	// by directly overwriting fields. Slices and maps are overwritten and not combined.
	// Subsequent calls to this API will merge inputs with the current configuration,
	// not the default configuration.
	// When updating an existing configuration, the revision number in the updated configuration
	// must match the revision number of the existing configuration, otherwise a conflict
	// error will be returned. The timestamp field of the revision is ignored.
	//
	// Note: some fields may contain secrets. The placeholder value "***" can be used to
	// keep an existing secret when updating the cluster configuration.
	SetConfiguration(context.Context, *CapabilityBackendConfigSpec) (*emptypb.Empty, error)
	// Resets the configuration of the managed Cortex cluster to the current default configuration.
	//
	// The request may optionally contain a field mask to specify which fields should
	// be preserved. Furthermore, if a mask is set, the request may also contain a patch
	// object used to apply additional changes to the masked fields. These changes are
	// applied atomically at the time of reset. Fields present in the patch object, but
	// not in the mask, are ignored.
	//
	// For example, with the following message:
	//
	//	message Example {
	//	  optional int32 a = 1;
	//	  optional int32 b = 2;
	//	  optional int32 c = 3;
	//	}
	//
	// and current state:
	//
	//	active:  { a: 1, b: 2, c: 3 }
	//	default: { a: 4, b: 5, c: 6 }
	//
	// and reset request parameters:
	//
	//	{
	//	  mask:    { paths: [ "a", "b" ] }
	//	  patch:   { a: 100 }
	//	}
	//
	// The resulting active configuration will be:
	//
	//	active:  {
	//	  a: 100, // masked, set to 100 via patch
	//	  b: 2,   // masked, but not set in patch, so left unchanged
	//	  c: 6,   // not masked, reset to default
	//	}
	ResetConfiguration(context.Context, *ResetRequest) (*emptypb.Empty, error)
	// Gets the current status of the managed Cortex cluster.
	// The status includes the current install state, version, and metadata. If
	// the cluster is in the process of being reconfigured or uninstalled, it will
	// be reflected in the install state.
	// No guarantees are made about the contents of the metadata field; its
	// contents are strictly informational.
	Status(context.Context, *emptypb.Empty) (*driverutil.InstallStatus, error)
	// Installs the managed Cortex cluster.
	// The cluster will be installed using the current configuration, or the
	// default configuration if none is explicitly set.
	Install(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// Uninstalls the managed Cortex cluster.
	// Implementation details including error handling and system state requirements
	// are left to the cluster driver, and this API makes no guarantees about
	// the state of the cluster after the call completes (regardless of success).
	Uninstall(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// Returns a static list of presets that can be used as a base for configuring the managed Cortex cluster.
	// There are several ways to use the presets, depending
	// on the desired behavior:
	//  1. Set the default configuration to a preset spec, then use SetConfiguration
	//     to fill in any additional required fields (credentials, etc)
	//  2. Add the required fields to the default configuration, then use
	//     SetConfiguration with a preset spec.
	//  3. Leave the default configuration as-is, and use SetConfiguration with a
	//     preset spec plus the required fields.
	ListPresets(context.Context, *emptypb.Empty) (*PresetList, error)
	// Show what changes would be made to a configuration without saving them.
	// The request expects an action, target, and spec to be provided. These
	// correspond roughly to the other APIs in this service.
	//
	// Configuring DryRunRequest:
	//   - Use the Active target for the SetConfiguration API, and the Default target
	//     for the SetDefaultConfiguration API. Install and Uninstall actions do not
	//     require a target.
	//   - Only the Set action requires a spec to be provided.
	//
	// Notes:
	//   - When DryRun is used on Install or Uninstall requests, the response will
	//     contain modifications to the 'enabled' field only. This field is read-only
	//     in the Set* APIs.
	//   - To validate the current configuration but keep it unchanged, use the
	//     Set action with an empty spec.
	//   - Configurations returned by DryRun will always have an empty revision field.
	DryRun(context.Context, *DryRunRequest) (*DryRunResponse, error)
	// Get a list of all past revisions of the configuration.
	// Will return the history for either the active or default configuration
	// depending on the specified target.
	// The entries are ordered from oldest to newest, where the last entry is
	// the current configuration.
	ConfigurationHistory(context.Context, *driverutil.ConfigurationHistoryRequest) (*ConfigurationHistoryResponse, error)
	mustEmbedUnimplementedCortexOpsServer()
}

// UnimplementedCortexOpsServer must be embedded to have forward compatible implementations.
type UnimplementedCortexOpsServer struct {
}

func (UnimplementedCortexOpsServer) GetDefaultConfiguration(context.Context, *driverutil.GetRequest) (*CapabilityBackendConfigSpec, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultConfiguration not implemented")
}
func (UnimplementedCortexOpsServer) SetDefaultConfiguration(context.Context, *CapabilityBackendConfigSpec) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDefaultConfiguration not implemented")
}
func (UnimplementedCortexOpsServer) ResetDefaultConfiguration(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetDefaultConfiguration not implemented")
}
func (UnimplementedCortexOpsServer) GetConfiguration(context.Context, *driverutil.GetRequest) (*CapabilityBackendConfigSpec, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfiguration not implemented")
}
func (UnimplementedCortexOpsServer) SetConfiguration(context.Context, *CapabilityBackendConfigSpec) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetConfiguration not implemented")
}
func (UnimplementedCortexOpsServer) ResetConfiguration(context.Context, *ResetRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetConfiguration not implemented")
}
func (UnimplementedCortexOpsServer) Status(context.Context, *emptypb.Empty) (*driverutil.InstallStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedCortexOpsServer) Install(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Install not implemented")
}
func (UnimplementedCortexOpsServer) Uninstall(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Uninstall not implemented")
}
func (UnimplementedCortexOpsServer) ListPresets(context.Context, *emptypb.Empty) (*PresetList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPresets not implemented")
}
func (UnimplementedCortexOpsServer) DryRun(context.Context, *DryRunRequest) (*DryRunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DryRun not implemented")
}
func (UnimplementedCortexOpsServer) ConfigurationHistory(context.Context, *driverutil.ConfigurationHistoryRequest) (*ConfigurationHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigurationHistory not implemented")
}
func (UnimplementedCortexOpsServer) mustEmbedUnimplementedCortexOpsServer() {}

// UnsafeCortexOpsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CortexOpsServer will
// result in compilation errors.
type UnsafeCortexOpsServer interface {
	mustEmbedUnimplementedCortexOpsServer()
}

func RegisterCortexOpsServer(s grpc.ServiceRegistrar, srv CortexOpsServer) {
	s.RegisterService(&CortexOps_ServiceDesc, srv)
}

func _CortexOps_GetDefaultConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(driverutil.GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).GetDefaultConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_GetDefaultConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).GetDefaultConfiguration(ctx, req.(*driverutil.GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_SetDefaultConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CapabilityBackendConfigSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).SetDefaultConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_SetDefaultConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).SetDefaultConfiguration(ctx, req.(*CapabilityBackendConfigSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_ResetDefaultConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).ResetDefaultConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_ResetDefaultConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).ResetDefaultConfiguration(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_GetConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(driverutil.GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).GetConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_GetConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).GetConfiguration(ctx, req.(*driverutil.GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_SetConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CapabilityBackendConfigSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).SetConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_SetConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).SetConfiguration(ctx, req.(*CapabilityBackendConfigSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_ResetConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).ResetConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_ResetConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).ResetConfiguration(ctx, req.(*ResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).Status(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_Install_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).Install(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_Install_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).Install(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_Uninstall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).Uninstall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_Uninstall_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).Uninstall(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_ListPresets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).ListPresets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_ListPresets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).ListPresets(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_DryRun_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DryRunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).DryRun(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_DryRun_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).DryRun(ctx, req.(*DryRunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CortexOps_ConfigurationHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(driverutil.ConfigurationHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CortexOpsServer).ConfigurationHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CortexOps_ConfigurationHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CortexOpsServer).ConfigurationHistory(ctx, req.(*driverutil.ConfigurationHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CortexOps_ServiceDesc is the grpc.ServiceDesc for CortexOps service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CortexOps_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cortexops.CortexOps",
	HandlerType: (*CortexOpsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDefaultConfiguration",
			Handler:    _CortexOps_GetDefaultConfiguration_Handler,
		},
		{
			MethodName: "SetDefaultConfiguration",
			Handler:    _CortexOps_SetDefaultConfiguration_Handler,
		},
		{
			MethodName: "ResetDefaultConfiguration",
			Handler:    _CortexOps_ResetDefaultConfiguration_Handler,
		},
		{
			MethodName: "GetConfiguration",
			Handler:    _CortexOps_GetConfiguration_Handler,
		},
		{
			MethodName: "SetConfiguration",
			Handler:    _CortexOps_SetConfiguration_Handler,
		},
		{
			MethodName: "ResetConfiguration",
			Handler:    _CortexOps_ResetConfiguration_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _CortexOps_Status_Handler,
		},
		{
			MethodName: "Install",
			Handler:    _CortexOps_Install_Handler,
		},
		{
			MethodName: "Uninstall",
			Handler:    _CortexOps_Uninstall_Handler,
		},
		{
			MethodName: "ListPresets",
			Handler:    _CortexOps_ListPresets_Handler,
		},
		{
			MethodName: "DryRun",
			Handler:    _CortexOps_DryRun_Handler,
		},
		{
			MethodName: "ConfigurationHistory",
			Handler:    _CortexOps_ConfigurationHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/aity-cloud/monty/plugins/metrics/apis/cortexops/cortexops.proto",
}
