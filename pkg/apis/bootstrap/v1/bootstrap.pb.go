// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v1.0.0
// source: github.com/aity-cloud/monty/pkg/apis/bootstrap/v1/bootstrap.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BootstrapJoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BootstrapJoinRequest) Reset() {
	*x = BootstrapJoinRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapJoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapJoinRequest) ProtoMessage() {}

func (x *BootstrapJoinRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapJoinRequest.ProtoReflect.Descriptor instead.
func (*BootstrapJoinRequest) Descriptor() ([]byte, []int) {
	return file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescGZIP(), []int{0}
}

type BootstrapJoinResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signatures map[string][]byte `protobuf:"bytes,1,rep,name=Signatures,proto3" json:"Signatures,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *BootstrapJoinResponse) Reset() {
	*x = BootstrapJoinResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapJoinResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapJoinResponse) ProtoMessage() {}

func (x *BootstrapJoinResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapJoinResponse.ProtoReflect.Descriptor instead.
func (*BootstrapJoinResponse) Descriptor() ([]byte, []int) {
	return file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescGZIP(), []int{1}
}

func (x *BootstrapJoinResponse) GetSignatures() map[string][]byte {
	if x != nil {
		return x.Signatures
	}
	return nil
}

type BootstrapAuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientID     string `protobuf:"bytes,1,opt,name=ClientID,proto3" json:"ClientID,omitempty"`
	ClientPubKey []byte `protobuf:"bytes,2,opt,name=ClientPubKey,proto3" json:"ClientPubKey,omitempty"`
	Capability   string `protobuf:"bytes,3,opt,name=Capability,proto3" json:"Capability,omitempty"`
}

func (x *BootstrapAuthRequest) Reset() {
	*x = BootstrapAuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapAuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapAuthRequest) ProtoMessage() {}

func (x *BootstrapAuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapAuthRequest.ProtoReflect.Descriptor instead.
func (*BootstrapAuthRequest) Descriptor() ([]byte, []int) {
	return file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescGZIP(), []int{2}
}

func (x *BootstrapAuthRequest) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *BootstrapAuthRequest) GetClientPubKey() []byte {
	if x != nil {
		return x.ClientPubKey
	}
	return nil
}

func (x *BootstrapAuthRequest) GetCapability() string {
	if x != nil {
		return x.Capability
	}
	return ""
}

type BootstrapAuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerPubKey []byte `protobuf:"bytes,1,opt,name=ServerPubKey,proto3" json:"ServerPubKey,omitempty"`
}

func (x *BootstrapAuthResponse) Reset() {
	*x = BootstrapAuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapAuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapAuthResponse) ProtoMessage() {}

func (x *BootstrapAuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapAuthResponse.ProtoReflect.Descriptor instead.
func (*BootstrapAuthResponse) Descriptor() ([]byte, []int) {
	return file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescGZIP(), []int{3}
}

func (x *BootstrapAuthResponse) GetServerPubKey() []byte {
	if x != nil {
		return x.ServerPubKey
	}
	return nil
}

var File_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto protoreflect.FileDescriptor

var file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDesc = []byte{
	0x0a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x74,
	0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70,
	0x2f, 0x76, 0x31, 0x2f, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x22, 0x16,
	0x0a, 0x14, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x4a, 0x6f, 0x69, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xa8, 0x01, 0x0a, 0x15, 0x42, 0x6f, 0x6f, 0x74, 0x73,
	0x74, 0x72, 0x61, 0x70, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x50, 0x0a, 0x0a, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70,
	0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x73, 0x1a, 0x3d, 0x0a, 0x0f, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x76, 0x0a, 0x14, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x41, 0x75,
	0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50,
	0x75, 0x62, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x61, 0x70,
	0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x43,
	0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x22, 0x3b, 0x0a, 0x15, 0x42, 0x6f, 0x6f,
	0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x75, 0x62, 0x4b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x32, 0xa1, 0x01, 0x0a, 0x09, 0x42, 0x6f, 0x6f, 0x74, 0x73,
	0x74, 0x72, 0x61, 0x70, 0x12, 0x49, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x1f, 0x2e, 0x62,
	0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72,
	0x61, 0x70, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e,
	0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74,
	0x72, 0x61, 0x70, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x49, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x1f, 0x2e, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74,
	0x72, 0x61, 0x70, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x41, 0x75, 0x74,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x62, 0x6f, 0x6f, 0x74, 0x73,
	0x74, 0x72, 0x61, 0x70, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x41, 0x75,
	0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x74, 0x79, 0x2d, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x2f, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescOnce sync.Once
	file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescData = file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDesc
)

func file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescGZIP() []byte {
	file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescOnce.Do(func() {
		file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescData)
	})
	return file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDescData
}

var file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_goTypes = []interface{}{
	(*BootstrapJoinRequest)(nil),  // 0: bootstrap.BootstrapJoinRequest
	(*BootstrapJoinResponse)(nil), // 1: bootstrap.BootstrapJoinResponse
	(*BootstrapAuthRequest)(nil),  // 2: bootstrap.BootstrapAuthRequest
	(*BootstrapAuthResponse)(nil), // 3: bootstrap.BootstrapAuthResponse
	nil,                           // 4: bootstrap.BootstrapJoinResponse.SignaturesEntry
}
var file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_depIdxs = []int32{
	4, // 0: bootstrap.BootstrapJoinResponse.Signatures:type_name -> bootstrap.BootstrapJoinResponse.SignaturesEntry
	0, // 1: bootstrap.Bootstrap.Join:input_type -> bootstrap.BootstrapJoinRequest
	2, // 2: bootstrap.Bootstrap.Auth:input_type -> bootstrap.BootstrapAuthRequest
	1, // 3: bootstrap.Bootstrap.Join:output_type -> bootstrap.BootstrapJoinResponse
	3, // 4: bootstrap.Bootstrap.Auth:output_type -> bootstrap.BootstrapAuthResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_init() }
func file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_init() {
	if File_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapJoinRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapJoinResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapAuthRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapAuthResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_goTypes,
		DependencyIndexes: file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_depIdxs,
		MessageInfos:      file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_msgTypes,
	}.Build()
	File_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto = out.File
	file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_rawDesc = nil
	file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_goTypes = nil
	file_github_com_aity_cloud_monty_pkg_apis_bootstrap_v1_bootstrap_proto_depIdxs = nil
}
