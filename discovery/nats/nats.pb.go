// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v4.25.1
// source: nats.proto

package nats

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

// MessageType defines the NATs message type
type MessageType int32

const (
	MessageType_MESSAGE_TYPE_REGISTER   MessageType = 0
	MessageType_MESSAGE_TYPE_DEREGISTER MessageType = 1
	MessageType_MESSAGE_TYPE_REQUEST    MessageType = 2
	MessageType_MESSAGE_TYPE_RESPONSE   MessageType = 3
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "MESSAGE_TYPE_REGISTER",
		1: "MESSAGE_TYPE_DEREGISTER",
		2: "MESSAGE_TYPE_REQUEST",
		3: "MESSAGE_TYPE_RESPONSE",
	}
	MessageType_value = map[string]int32{
		"MESSAGE_TYPE_REGISTER":   0,
		"MESSAGE_TYPE_DEREGISTER": 1,
		"MESSAGE_TYPE_REQUEST":    2,
		"MESSAGE_TYPE_RESPONSE":   3,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_nats_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_nats_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_nats_proto_rawDescGZIP(), []int{0}
}

// Message defines the NATs message
// used by the discovery provider
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Specifies the host name of the client node
	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	// Specifies the port of the client node
	Port int32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	// Specifies the client name
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// Specifies the message type
	MessageType MessageType `protobuf:"varint,4,opt,name=message_type,json=messageType,proto3,enum=nats.MessageType" json:"message_type,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_nats_proto_msgTypes[0]
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
	return file_nats_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *Message) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Message) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Message) GetMessageType() MessageType {
	if x != nil {
		return x.MessageType
	}
	return MessageType_MESSAGE_TYPE_REGISTER
}

var File_nats_proto protoreflect.FileDescriptor

var file_nats_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6e, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6e, 0x61,
	0x74, 0x73, 0x22, 0x7b, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x0c, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x6e, 0x61, 0x74, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x2a,
	0x7a, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19,
	0x0a, 0x15, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52,
	0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x10, 0x00, 0x12, 0x1b, 0x0a, 0x17, 0x4d, 0x45, 0x53,
	0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x45, 0x52, 0x45, 0x47, 0x49,
	0x53, 0x54, 0x45, 0x52, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x02,
	0x12, 0x19, 0x0a, 0x15, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x03, 0x42, 0x08, 0x5a, 0x06, 0x2e,
	0x3b, 0x6e, 0x61, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nats_proto_rawDescOnce sync.Once
	file_nats_proto_rawDescData = file_nats_proto_rawDesc
)

func file_nats_proto_rawDescGZIP() []byte {
	file_nats_proto_rawDescOnce.Do(func() {
		file_nats_proto_rawDescData = protoimpl.X.CompressGZIP(file_nats_proto_rawDescData)
	})
	return file_nats_proto_rawDescData
}

var file_nats_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_nats_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_nats_proto_goTypes = []interface{}{
	(MessageType)(0), // 0: nats.MessageType
	(*Message)(nil),  // 1: nats.Message
}
var file_nats_proto_depIdxs = []int32{
	0, // 0: nats.Message.message_type:type_name -> nats.MessageType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_nats_proto_init() }
func file_nats_proto_init() {
	if File_nats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nats_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_nats_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_nats_proto_goTypes,
		DependencyIndexes: file_nats_proto_depIdxs,
		EnumInfos:         file_nats_proto_enumTypes,
		MessageInfos:      file_nats_proto_msgTypes,
	}.Build()
	File_nats_proto = out.File
	file_nats_proto_rawDesc = nil
	file_nats_proto_goTypes = nil
	file_nats_proto_depIdxs = nil
}
