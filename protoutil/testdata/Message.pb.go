// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: Message.proto

package testdata

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_	= protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_	= protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Message struct {
	state		protoimpl.MessageState
	sizeCache	protoimpl.SizeCache
	unknownFields	protoimpl.UnknownFields

	Uid		string			`protobuf:"bytes,10,opt,name=uid,proto3" json:"uid,omitempty"`
	TimestampUtc	*timestamp.Timestamp	`protobuf:"bytes,11,opt,name=timestamp_utc,json=timestampUtc,proto3" json:"timestamp_utc,omitempty"`
	Elapsed		*duration.Duration	`protobuf:"bytes,12,opt,name=elapsed,proto3" json:"elapsed,omitempty"`
	StatusCode	int32			`protobuf:"varint,13,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	ContentLength	int64			`protobuf:"varint,14,opt,name=content_length,json=contentLength,proto3" json:"content_length,omitempty"`
	Value		float64			`protobuf:"fixed64,15,opt,name=value,proto3" json:"value,omitempty"`
	Error		string			`protobuf:"bytes,16,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage()	{}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_Message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_Message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *Message) GetTimestampUtc() *timestamp.Timestamp {
	if x != nil {
		return x.TimestampUtc
	}
	return nil
}

func (x *Message) GetElapsed() *duration.Duration {
	if x != nil {
		return x.Elapsed
	}
	return nil
}

func (x *Message) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Message) GetContentLength() int64 {
	if x != nil {
		return x.ContentLength
	}
	return 0
}

func (x *Message) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Message) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_Message_proto protoreflect.FileDescriptor

var file_Message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x02, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x3f, 0x0a, 0x0d, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x5f, 0x75, 0x74, 0x63, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x55, 0x74, 0x63, 0x12, 0x33, 0x0a, 0x07, 0x65, 0x6c, 0x61,
	0x70, 0x73, 0x65, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Message_proto_rawDescOnce	sync.Once
	file_Message_proto_rawDescData	= file_Message_proto_rawDesc
)

func file_Message_proto_rawDescGZIP() []byte {
	file_Message_proto_rawDescOnce.Do(func() {
		file_Message_proto_rawDescData = protoimpl.X.CompressGZIP(file_Message_proto_rawDescData)
	})
	return file_Message_proto_rawDescData
}

var file_Message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_Message_proto_goTypes = []interface{}{
	(*Message)(nil),		// 0: testdata.Message
	(*timestamp.Timestamp)(nil),	// 1: google.protobuf.Timestamp
	(*duration.Duration)(nil),	// 2: google.protobuf.Duration
}
var file_Message_proto_depIdxs = []int32{
	1,	// 0: testdata.Message.timestamp_utc:type_name -> google.protobuf.Timestamp
	2,	// 1: testdata.Message.elapsed:type_name -> google.protobuf.Duration
	2,	// [2:2] is the sub-list for method output_type
	2,	// [2:2] is the sub-list for method input_type
	2,	// [2:2] is the sub-list for extension type_name
	2,	// [2:2] is the sub-list for extension extendee
	0,	// [0:2] is the sub-list for field type_name
}

func init()	{ file_Message_proto_init() }
func file_Message_proto_init() {
	if File_Message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			GoPackagePath:	reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor:	file_Message_proto_rawDesc,
			NumEnums:	0,
			NumMessages:	1,
			NumExtensions:	0,
			NumServices:	0,
		},
		GoTypes:		file_Message_proto_goTypes,
		DependencyIndexes:	file_Message_proto_depIdxs,
		MessageInfos:		file_Message_proto_msgTypes,
	}.Build()
	File_Message_proto = out.File
	file_Message_proto_rawDesc = nil
	file_Message_proto_goTypes = nil
	file_Message_proto_depIdxs = nil
}
