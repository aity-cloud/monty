// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v1.0.0
// source: github.com/aity-cloud/monty/pkg/test/testgrpc/stream.proto

package testgrpc

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

type StreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Request string `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescGZIP(), []int{0}
}

func (x *StreamRequest) GetRequest() string {
	if x != nil {
		return x.Request
	}
	return ""
}

type StreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *StreamResponse) Reset() {
	*x = StreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse) ProtoMessage() {}

func (x *StreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResponse.ProtoReflect.Descriptor instead.
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescGZIP(), []int{1}
}

func (x *StreamResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

var File_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto protoreflect.FileDescriptor

var file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDesc = []byte{
	0x0a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x74,
	0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2f,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x74, 0x65,
	0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x29, 0x0a,
	0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x2c, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x5e, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x12, 0x1e, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x74, 0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f,
	0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x74,
	0x65, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescOnce sync.Once
	file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescData = file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDesc
)

func file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescGZIP() []byte {
	file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescOnce.Do(func() {
		file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescData)
	})
	return file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDescData
}

var file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_goTypes = []interface{}{
	(*StreamRequest)(nil),  // 0: testgrpc.stream.StreamRequest
	(*StreamResponse)(nil), // 1: testgrpc.stream.StreamResponse
}
var file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_depIdxs = []int32{
	0, // 0: testgrpc.stream.StreamService.Stream:input_type -> testgrpc.stream.StreamRequest
	1, // 1: testgrpc.stream.StreamService.Stream:output_type -> testgrpc.stream.StreamResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_init() }
func file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_init() {
	if File_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamRequest); i {
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
		file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamResponse); i {
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
			RawDescriptor: file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_goTypes,
		DependencyIndexes: file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_depIdxs,
		MessageInfos:      file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_msgTypes,
	}.Build()
	File_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto = out.File
	file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_rawDesc = nil
	file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_goTypes = nil
	file_github_com_aity_cloud_monty_pkg_test_testgrpc_stream_proto_depIdxs = nil
}
