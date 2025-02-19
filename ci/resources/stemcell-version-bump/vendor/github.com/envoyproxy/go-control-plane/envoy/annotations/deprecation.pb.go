// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v5.26.1
// source: envoy/annotations/deprecation.proto

package annotations

import (
	reflect "reflect"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_envoy_annotations_deprecation_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         189503207,
		Name:          "envoy.annotations.disallowed_by_default",
		Tag:           "varint,189503207,opt,name=disallowed_by_default",
		Filename:      "envoy/annotations/deprecation.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         157299826,
		Name:          "envoy.annotations.deprecated_at_minor_version",
		Tag:           "bytes,157299826,opt,name=deprecated_at_minor_version",
		Filename:      "envoy/annotations/deprecation.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumValueOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         70100853,
		Name:          "envoy.annotations.disallowed_by_default_enum",
		Tag:           "varint,70100853,opt,name=disallowed_by_default_enum",
		Filename:      "envoy/annotations/deprecation.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumValueOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         181198657,
		Name:          "envoy.annotations.deprecated_at_minor_version_enum",
		Tag:           "bytes,181198657,opt,name=deprecated_at_minor_version_enum",
		Filename:      "envoy/annotations/deprecation.proto",
	},
}

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional bool disallowed_by_default = 189503207;
	E_DisallowedByDefault = &file_envoy_annotations_deprecation_proto_extTypes[0]
	// The API major and minor version on which the field was deprecated
	// (e.g., "3.5" for major version 3 and minor version 5).
	//
	// optional string deprecated_at_minor_version = 157299826;
	E_DeprecatedAtMinorVersion = &file_envoy_annotations_deprecation_proto_extTypes[1]
)

// Extension fields to descriptorpb.EnumValueOptions.
var (
	// optional bool disallowed_by_default_enum = 70100853;
	E_DisallowedByDefaultEnum = &file_envoy_annotations_deprecation_proto_extTypes[2]
	// The API major and minor version on which the enum value was deprecated
	// (e.g., "3.5" for major version 3 and minor version 5).
	//
	// optional string deprecated_at_minor_version_enum = 181198657;
	E_DeprecatedAtMinorVersionEnum = &file_envoy_annotations_deprecation_proto_extTypes[3]
)

var File_envoy_annotations_deprecation_proto protoreflect.FileDescriptor

var file_envoy_annotations_deprecation_proto_rawDesc = []byte{
	0x0a, 0x23, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2f, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x54, 0x0a, 0x15, 0x64, 0x69,
	0x73, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x5f, 0x64, 0x65, 0x66, 0x61,
	0x75, 0x6c, 0x74, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xe7, 0xad, 0xae, 0x5a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x64, 0x69, 0x73,
	0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x42, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x3a, 0x5f, 0x0a, 0x1b, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x5f, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf2,
	0xe8, 0x80, 0x4b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x18, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x4d, 0x69, 0x6e, 0x6f, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x3a, 0x61, 0x0a, 0x1a, 0x64, 0x69, 0x73, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f,
	0x62, 0x79, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x12,
	0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xf5, 0xce, 0xb6, 0x21, 0x20, 0x01, 0x28, 0x08, 0x52, 0x17, 0x64, 0x69, 0x73,
	0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x42, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x45, 0x6e, 0x75, 0x6d, 0x3a, 0x6c, 0x0a, 0x20, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x5f, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x12, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc1, 0xbe, 0xb3, 0x56,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x1c, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x4d, 0x69, 0x6e, 0x6f, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e,
	0x75, 0x6d, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_envoy_annotations_deprecation_proto_goTypes = []interface{}{
	(*descriptorpb.FieldOptions)(nil),     // 0: google.protobuf.FieldOptions
	(*descriptorpb.EnumValueOptions)(nil), // 1: google.protobuf.EnumValueOptions
}
var file_envoy_annotations_deprecation_proto_depIdxs = []int32{
	0, // 0: envoy.annotations.disallowed_by_default:extendee -> google.protobuf.FieldOptions
	0, // 1: envoy.annotations.deprecated_at_minor_version:extendee -> google.protobuf.FieldOptions
	1, // 2: envoy.annotations.disallowed_by_default_enum:extendee -> google.protobuf.EnumValueOptions
	1, // 3: envoy.annotations.deprecated_at_minor_version_enum:extendee -> google.protobuf.EnumValueOptions
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	0, // [0:4] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_envoy_annotations_deprecation_proto_init() }
func file_envoy_annotations_deprecation_proto_init() {
	if File_envoy_annotations_deprecation_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_envoy_annotations_deprecation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_envoy_annotations_deprecation_proto_goTypes,
		DependencyIndexes: file_envoy_annotations_deprecation_proto_depIdxs,
		ExtensionInfos:    file_envoy_annotations_deprecation_proto_extTypes,
	}.Build()
	File_envoy_annotations_deprecation_proto = out.File
	file_envoy_annotations_deprecation_proto_rawDesc = nil
	file_envoy_annotations_deprecation_proto_goTypes = nil
	file_envoy_annotations_deprecation_proto_depIdxs = nil
}
