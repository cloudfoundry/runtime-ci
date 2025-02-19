// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.0--rc2
// source: xds/core/v3/extension.proto

package v3

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TypedExtensionConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	TypedConfig *anypb.Any `protobuf:"bytes,2,opt,name=typed_config,json=typedConfig,proto3" json:"typed_config,omitempty"`
}

func (x *TypedExtensionConfig) Reset() {
	*x = TypedExtensionConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_core_v3_extension_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TypedExtensionConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypedExtensionConfig) ProtoMessage() {}

func (x *TypedExtensionConfig) ProtoReflect() protoreflect.Message {
	mi := &file_xds_core_v3_extension_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypedExtensionConfig.ProtoReflect.Descriptor instead.
func (*TypedExtensionConfig) Descriptor() ([]byte, []int) {
	return file_xds_core_v3_extension_proto_rawDescGZIP(), []int{0}
}

func (x *TypedExtensionConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TypedExtensionConfig) GetTypedConfig() *anypb.Any {
	if x != nil {
		return x.TypedConfig
	}
	return nil
}

var File_xds_core_v3_extension_proto protoreflect.FileDescriptor

var file_xds_core_v3_extension_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x78, 0x64, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x33, 0x2f, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x78,
	0x64, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x33, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76,
	0x0a, 0x14, 0x54, 0x79, 0x70, 0x65, 0x64, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x41, 0x0a, 0x0c, 0x74, 0x79, 0x70, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x42,
	0x08, 0xfa, 0x42, 0x05, 0xa2, 0x01, 0x02, 0x08, 0x01, 0x52, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x64,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x4e, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x78, 0x64, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x33,
	0x42, 0x0e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63,
	0x6e, 0x63, 0x66, 0x2f, 0x78, 0x64, 0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x78, 0x64, 0x73, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2f, 0x76, 0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_xds_core_v3_extension_proto_rawDescOnce sync.Once
	file_xds_core_v3_extension_proto_rawDescData = file_xds_core_v3_extension_proto_rawDesc
)

func file_xds_core_v3_extension_proto_rawDescGZIP() []byte {
	file_xds_core_v3_extension_proto_rawDescOnce.Do(func() {
		file_xds_core_v3_extension_proto_rawDescData = protoimpl.X.CompressGZIP(file_xds_core_v3_extension_proto_rawDescData)
	})
	return file_xds_core_v3_extension_proto_rawDescData
}

var file_xds_core_v3_extension_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_xds_core_v3_extension_proto_goTypes = []interface{}{
	(*TypedExtensionConfig)(nil), // 0: xds.core.v3.TypedExtensionConfig
	(*anypb.Any)(nil),            // 1: google.protobuf.Any
}
var file_xds_core_v3_extension_proto_depIdxs = []int32{
	1, // 0: xds.core.v3.TypedExtensionConfig.typed_config:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_xds_core_v3_extension_proto_init() }
func file_xds_core_v3_extension_proto_init() {
	if File_xds_core_v3_extension_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_xds_core_v3_extension_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TypedExtensionConfig); i {
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
			RawDescriptor: file_xds_core_v3_extension_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_xds_core_v3_extension_proto_goTypes,
		DependencyIndexes: file_xds_core_v3_extension_proto_depIdxs,
		MessageInfos:      file_xds_core_v3_extension_proto_msgTypes,
	}.Build()
	File_xds_core_v3_extension_proto = out.File
	file_xds_core_v3_extension_proto_rawDesc = nil
	file_xds_core_v3_extension_proto_goTypes = nil
	file_xds_core_v3_extension_proto_depIdxs = nil
}
