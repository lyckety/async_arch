// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: eventheaders/header/v1/header.proto

package headerv1

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

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId       string `protobuf:"bytes,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	EventVersion  string `protobuf:"bytes,2,opt,name=event_version,json=eventVersion,proto3" json:"event_version,omitempty"`
	EventName     string `protobuf:"bytes,3,opt,name=event_name,json=eventName,proto3" json:"event_name,omitempty"`
	EventTime     int64  `protobuf:"varint,4,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	EventProducer string `protobuf:"bytes,5,opt,name=event_producer,json=eventProducer,proto3" json:"event_producer,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventheaders_header_v1_header_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_eventheaders_header_v1_header_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_eventheaders_header_v1_header_proto_rawDescGZIP(), []int{0}
}

func (x *Header) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *Header) GetEventVersion() string {
	if x != nil {
		return x.EventVersion
	}
	return ""
}

func (x *Header) GetEventName() string {
	if x != nil {
		return x.EventName
	}
	return ""
}

func (x *Header) GetEventTime() int64 {
	if x != nil {
		return x.EventTime
	}
	return 0
}

func (x *Header) GetEventProducer() string {
	if x != nil {
		return x.EventProducer
	}
	return ""
}

var File_eventheaders_header_v1_header_proto protoreflect.FileDescriptor

var file_eventheaders_header_v1_header_proto_rawDesc = []byte{
	0x0a, 0x23, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x2f, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x2e, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0xad, 0x01,
	0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x42, 0x84, 0x02,
	0x0a, 0x1a, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x2e, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x79, 0x63, 0x6b, 0x65, 0x74, 0x79, 0x2f,
	0x61, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x70, 0x6f, 0x70, 0x75, 0x67,
	0x5f, 0x6a, 0x69, 0x72, 0x61, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2d, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x2f, 0x76, 0x31, 0x3b, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x45,
	0x48, 0x58, 0xaa, 0x02, 0x16, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x16, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x5c, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x22, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x5c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x18, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x3a, 0x3a, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_eventheaders_header_v1_header_proto_rawDescOnce sync.Once
	file_eventheaders_header_v1_header_proto_rawDescData = file_eventheaders_header_v1_header_proto_rawDesc
)

func file_eventheaders_header_v1_header_proto_rawDescGZIP() []byte {
	file_eventheaders_header_v1_header_proto_rawDescOnce.Do(func() {
		file_eventheaders_header_v1_header_proto_rawDescData = protoimpl.X.CompressGZIP(file_eventheaders_header_v1_header_proto_rawDescData)
	})
	return file_eventheaders_header_v1_header_proto_rawDescData
}

var file_eventheaders_header_v1_header_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_eventheaders_header_v1_header_proto_goTypes = []interface{}{
	(*Header)(nil), // 0: eventheaders.header.v1.Header
}
var file_eventheaders_header_v1_header_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_eventheaders_header_v1_header_proto_init() }
func file_eventheaders_header_v1_header_proto_init() {
	if File_eventheaders_header_v1_header_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_eventheaders_header_v1_header_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
			RawDescriptor: file_eventheaders_header_v1_header_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_eventheaders_header_v1_header_proto_goTypes,
		DependencyIndexes: file_eventheaders_header_v1_header_proto_depIdxs,
		MessageInfos:      file_eventheaders_header_v1_header_proto_msgTypes,
	}.Build()
	File_eventheaders_header_v1_header_proto = out.File
	file_eventheaders_header_v1_header_proto_rawDesc = nil
	file_eventheaders_header_v1_header_proto_goTypes = nil
	file_eventheaders_header_v1_header_proto_depIdxs = nil
}