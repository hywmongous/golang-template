// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: identity_snapshots.proto

package cqrs

import (
	reflect "reflect"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_identity_snapshots_proto protoreflect.FileDescriptor

var file_identity_snapshots_proto_rawDesc = []byte{
	0x0a, 0x18, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x71, 0x72, 0x73,
	0x42, 0x21, 0x5a, 0x1f, 0x2e, 0x2e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x69, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63,
	0x71, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_identity_snapshots_proto_goTypes = []interface{}{}
	file_identity_snapshots_proto_depIdxs = []int32{
		0, // [0:0] is the sub-list for method output_type
		0, // [0:0] is the sub-list for method input_type
		0, // [0:0] is the sub-list for extension type_name
		0, // [0:0] is the sub-list for extension extendee
		0, // [0:0] is the sub-list for field type_name
	}
)

func init() { file_identity_snapshots_proto_init() }
func file_identity_snapshots_proto_init() {
	if File_identity_snapshots_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_identity_snapshots_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_identity_snapshots_proto_goTypes,
		DependencyIndexes: file_identity_snapshots_proto_depIdxs,
	}.Build()
	File_identity_snapshots_proto = out.File
	file_identity_snapshots_proto_rawDesc = nil
	file_identity_snapshots_proto_goTypes = nil
	file_identity_snapshots_proto_depIdxs = nil
}
