// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: kafkaconsumer/v1/kafkaconsumer.proto

package kafkaconsumer

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

type SubscribeTopicRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *SubscribeTopicRequest) Reset() {
	*x = SubscribeTopicRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeTopicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeTopicRequest) ProtoMessage() {}

func (x *SubscribeTopicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeTopicRequest.ProtoReflect.Descriptor instead.
func (*SubscribeTopicRequest) Descriptor() ([]byte, []int) {
	return file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescGZIP(), []int{0}
}

func (x *SubscribeTopicRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type SubscribeTopicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key       string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value     string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Topic     string `protobuf:"bytes,3,opt,name=topic,proto3" json:"topic,omitempty"`
	Timestamp uint64 `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *SubscribeTopicResponse) Reset() {
	*x = SubscribeTopicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeTopicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeTopicResponse) ProtoMessage() {}

func (x *SubscribeTopicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeTopicResponse.ProtoReflect.Descriptor instead.
func (*SubscribeTopicResponse) Descriptor() ([]byte, []int) {
	return file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescGZIP(), []int{1}
}

func (x *SubscribeTopicResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SubscribeTopicResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *SubscribeTopicResponse) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *SubscribeTopicResponse) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_kafkaconsumer_v1_kafkaconsumer_proto protoreflect.FileDescriptor

var file_kafkaconsumer_v1_kafkaconsumer_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x2f,
	0x76, 0x31, 0x2f, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x63, 0x6f, 0x6e,
	0x73, 0x75, 0x6d, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x2d, 0x0a, 0x15, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x74, 0x0a, 0x16, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x32, 0x7d, 0x0a,
	0x14, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x65, 0x0a, 0x0e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x27, 0x2e, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x63,
	0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x28, 0x2e, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x70,
	0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x56, 0x5a, 0x54,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x79, 0x63, 0x6b, 0x65,
	0x74, 0x79, 0x2f, 0x61, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x67, 0x6f,
	0x6c, 0x61, 0x6e, 0x67, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x63, 0x6f, 0x6e, 0x73, 0x75,
	0x6d, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x63, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescOnce sync.Once
	file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescData = file_kafkaconsumer_v1_kafkaconsumer_proto_rawDesc
)

func file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescGZIP() []byte {
	file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescOnce.Do(func() {
		file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescData = protoimpl.X.CompressGZIP(file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescData)
	})
	return file_kafkaconsumer_v1_kafkaconsumer_proto_rawDescData
}

var file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kafkaconsumer_v1_kafkaconsumer_proto_goTypes = []interface{}{
	(*SubscribeTopicRequest)(nil),  // 0: kafkaconsumer.v1.SubscribeTopicRequest
	(*SubscribeTopicResponse)(nil), // 1: kafkaconsumer.v1.SubscribeTopicResponse
}
var file_kafkaconsumer_v1_kafkaconsumer_proto_depIdxs = []int32{
	0, // 0: kafkaconsumer.v1.KafkaConsumerService.SubscribeTopic:input_type -> kafkaconsumer.v1.SubscribeTopicRequest
	1, // 1: kafkaconsumer.v1.KafkaConsumerService.SubscribeTopic:output_type -> kafkaconsumer.v1.SubscribeTopicResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_kafkaconsumer_v1_kafkaconsumer_proto_init() }
func file_kafkaconsumer_v1_kafkaconsumer_proto_init() {
	if File_kafkaconsumer_v1_kafkaconsumer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeTopicRequest); i {
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
		file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeTopicResponse); i {
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
			RawDescriptor: file_kafkaconsumer_v1_kafkaconsumer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kafkaconsumer_v1_kafkaconsumer_proto_goTypes,
		DependencyIndexes: file_kafkaconsumer_v1_kafkaconsumer_proto_depIdxs,
		MessageInfos:      file_kafkaconsumer_v1_kafkaconsumer_proto_msgTypes,
	}.Build()
	File_kafkaconsumer_v1_kafkaconsumer_proto = out.File
	file_kafkaconsumer_v1_kafkaconsumer_proto_rawDesc = nil
	file_kafkaconsumer_v1_kafkaconsumer_proto_goTypes = nil
	file_kafkaconsumer_v1_kafkaconsumer_proto_depIdxs = nil
}
