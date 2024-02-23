// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: kafkaproducer/v1/kafkaproducer.proto

package kafkaproducer

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

type SendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key       string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value     string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Topic     string `protobuf:"bytes,3,opt,name=topic,proto3" json:"topic,omitempty"`
	Timestamp uint64 `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *SendRequest) Reset() {
	*x = SendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kafkaproducer_v1_kafkaproducer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRequest) ProtoMessage() {}

func (x *SendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kafkaproducer_v1_kafkaproducer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRequest.ProtoReflect.Descriptor instead.
func (*SendRequest) Descriptor() ([]byte, []int) {
	return file_kafkaproducer_v1_kafkaproducer_proto_rawDescGZIP(), []int{0}
}

func (x *SendRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SendRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *SendRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *SendRequest) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type SendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendResponse) Reset() {
	*x = SendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kafkaproducer_v1_kafkaproducer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendResponse) ProtoMessage() {}

func (x *SendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kafkaproducer_v1_kafkaproducer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendResponse.ProtoReflect.Descriptor instead.
func (*SendResponse) Descriptor() ([]byte, []int) {
	return file_kafkaproducer_v1_kafkaproducer_proto_rawDescGZIP(), []int{1}
}

var File_kafkaproducer_v1_kafkaproducer_proto protoreflect.FileDescriptor

var file_kafkaproducer_v1_kafkaproducer_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x2f,
	0x76, 0x31, 0x2f, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x69, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x22, 0x0e, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0x5f, 0x0a, 0x14, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x04, 0x53,
	0x65, 0x6e, 0x64, 0x12, 0x1d, 0x2e, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x56, 0x5a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6c, 0x79, 0x63, 0x6b, 0x65, 0x74, 0x79, 0x2f, 0x61, 0x73, 0x79, 0x6e, 0x63,
	0x5f, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6b, 0x61,
	0x66, 0x6b, 0x61, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6b,
	0x61, 0x66, 0x6b, 0x61, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kafkaproducer_v1_kafkaproducer_proto_rawDescOnce sync.Once
	file_kafkaproducer_v1_kafkaproducer_proto_rawDescData = file_kafkaproducer_v1_kafkaproducer_proto_rawDesc
)

func file_kafkaproducer_v1_kafkaproducer_proto_rawDescGZIP() []byte {
	file_kafkaproducer_v1_kafkaproducer_proto_rawDescOnce.Do(func() {
		file_kafkaproducer_v1_kafkaproducer_proto_rawDescData = protoimpl.X.CompressGZIP(file_kafkaproducer_v1_kafkaproducer_proto_rawDescData)
	})
	return file_kafkaproducer_v1_kafkaproducer_proto_rawDescData
}

var file_kafkaproducer_v1_kafkaproducer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kafkaproducer_v1_kafkaproducer_proto_goTypes = []interface{}{
	(*SendRequest)(nil),  // 0: kafkaproducer.v1.SendRequest
	(*SendResponse)(nil), // 1: kafkaproducer.v1.SendResponse
}
var file_kafkaproducer_v1_kafkaproducer_proto_depIdxs = []int32{
	0, // 0: kafkaproducer.v1.KafkaProducerService.Send:input_type -> kafkaproducer.v1.SendRequest
	1, // 1: kafkaproducer.v1.KafkaProducerService.Send:output_type -> kafkaproducer.v1.SendResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_kafkaproducer_v1_kafkaproducer_proto_init() }
func file_kafkaproducer_v1_kafkaproducer_proto_init() {
	if File_kafkaproducer_v1_kafkaproducer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kafkaproducer_v1_kafkaproducer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRequest); i {
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
		file_kafkaproducer_v1_kafkaproducer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendResponse); i {
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
			RawDescriptor: file_kafkaproducer_v1_kafkaproducer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kafkaproducer_v1_kafkaproducer_proto_goTypes,
		DependencyIndexes: file_kafkaproducer_v1_kafkaproducer_proto_depIdxs,
		MessageInfos:      file_kafkaproducer_v1_kafkaproducer_proto_msgTypes,
	}.Build()
	File_kafkaproducer_v1_kafkaproducer_proto = out.File
	file_kafkaproducer_v1_kafkaproducer_proto_rawDesc = nil
	file_kafkaproducer_v1_kafkaproducer_proto_goTypes = nil
	file_kafkaproducer_v1_kafkaproducer_proto_depIdxs = nil
}