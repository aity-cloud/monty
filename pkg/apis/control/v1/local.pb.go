// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v1.0.0
// source: github.com/aity-cloud/monty/pkg/apis/control/v1/local.proto

package v1

import (
	v1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	_ "github.com/kralicky/totem"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto protoreflect.FileDescriptor

var file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_rawDesc = []byte{
	0x0a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x74,
	0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x76,
	0x31, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x69, 0x74, 0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x72, 0x61, 0x6c, 0x69, 0x63, 0x6b,
	0x79, 0x2f, 0x74, 0x6f, 0x74, 0x65, 0x6d, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x45, 0x0a, 0x08, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x12, 0x31, 0x0a, 0x06, 0x57, 0x68, 0x6f, 0x61, 0x6d, 0x69, 0x12, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x1a, 0x06, 0x92, 0xf1, 0x04, 0x02, 0x08, 0x01, 0x42,
	0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69,
	0x74, 0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_goTypes = []interface{}{
	(*emptypb.Empty)(nil), // 0: google.protobuf.Empty
	(*v1.Reference)(nil),  // 1: core.Reference
}
var file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_depIdxs = []int32{
	0, // 0: control.Identity.Whoami:input_type -> google.protobuf.Empty
	1, // 1: control.Identity.Whoami:output_type -> core.Reference
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_init() }
func file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_init() {
	if File_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_goTypes,
		DependencyIndexes: file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_depIdxs,
	}.Build()
	File_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto = out.File
	file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_rawDesc = nil
	file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_goTypes = nil
	file_github_com_aity_cloud_monty_pkg_apis_control_v1_local_proto_depIdxs = nil
}
