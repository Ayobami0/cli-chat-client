// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.3
// source: message_message.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message_MessageType int32

const (
	Message_MESSAGE_TYPE_UNSPECIFIED  Message_MessageType = 0
	Message_MESSAGE_TYPE_NOTIFICATION Message_MessageType = 1
	Message_MESSAGE_TYPE_REGULAR      Message_MessageType = 2
)

// Enum value maps for Message_MessageType.
var (
	Message_MessageType_name = map[int32]string{
		0: "MESSAGE_TYPE_UNSPECIFIED",
		1: "MESSAGE_TYPE_NOTIFICATION",
		2: "MESSAGE_TYPE_REGULAR",
	}
	Message_MessageType_value = map[string]int32{
		"MESSAGE_TYPE_UNSPECIFIED":  0,
		"MESSAGE_TYPE_NOTIFICATION": 1,
		"MESSAGE_TYPE_REGULAR":      2,
	}
)

func (x Message_MessageType) Enum() *Message_MessageType {
	p := new(Message_MessageType)
	*p = x
	return p
}

func (x Message_MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Message_MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_message_message_proto_enumTypes[0].Descriptor()
}

func (Message_MessageType) Type() protoreflect.EnumType {
	return &file_message_message_proto_enumTypes[0]
}

func (x Message_MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Message_MessageType.Descriptor instead.
func (Message_MessageType) EnumDescriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{0, 0}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Sender  *User                  `protobuf:"bytes,2,opt,name=sender,proto3,oneof" json:"sender,omitempty"`
	Type    Message_MessageType    `protobuf:"varint,3,opt,name=type,proto3,enum=chat.Message_MessageType" json:"type,omitempty"`
	Content string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	SentAt  *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=sent_at,json=sentAt,proto3" json:"sent_at,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[0]
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
	return file_message_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Message) GetSender() *User {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *Message) GetType() Message_MessageType {
	if x != nil {
		return x.Type
	}
	return Message_MESSAGE_TYPE_UNSPECIFIED
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetSentAt() *timestamppb.Timestamp {
	if x != nil {
		return x.SentAt
	}
	return nil
}

type MessageStream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatId  string   `protobuf:"bytes,1,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
	Message *Message `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MessageStream) Reset() {
	*x = MessageStream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageStream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageStream) ProtoMessage() {}

func (x *MessageStream) ProtoReflect() protoreflect.Message {
	mi := &file_message_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageStream.ProtoReflect.Descriptor instead.
func (*MessageStream) Descriptor() ([]byte, []int) {
	return file_message_message_proto_rawDescGZIP(), []int{1}
}

func (x *MessageStream) GetChatId() string {
	if x != nil {
		return x.ChatId
	}
	return ""
}

func (x *MessageStream) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_message_message_proto protoreflect.FileDescriptor

var file_message_message_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x68, 0x61, 0x74, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xb1, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27,
	0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x48, 0x00, 0x52, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x33, 0x0a, 0x07, 0x73, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x73,
	0x65, 0x6e, 0x74, 0x41, 0x74, 0x22, 0x64, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x18, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x1d, 0x0a, 0x19, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10,
	0x01, 0x12, 0x18, 0x0a, 0x14, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x52, 0x45, 0x47, 0x55, 0x4c, 0x41, 0x52, 0x10, 0x02, 0x42, 0x09, 0x0a, 0x07, 0x5f,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x22, 0x51, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64,
	0x12, 0x27, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x79, 0x6f, 0x62, 0x61, 0x6d, 0x69, 0x30,
	0x2f, 0x63, 0x6c, 0x69, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_message_proto_rawDescOnce sync.Once
	file_message_message_proto_rawDescData = file_message_message_proto_rawDesc
)

func file_message_message_proto_rawDescGZIP() []byte {
	file_message_message_proto_rawDescOnce.Do(func() {
		file_message_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_message_proto_rawDescData)
	})
	return file_message_message_proto_rawDescData
}

var file_message_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_message_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_message_proto_goTypes = []interface{}{
	(Message_MessageType)(0),      // 0: chat.Message.MessageType
	(*Message)(nil),               // 1: chat.Message
	(*MessageStream)(nil),         // 2: chat.MessageStream
	(*User)(nil),                  // 3: chat.User
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_message_message_proto_depIdxs = []int32{
	3, // 0: chat.Message.sender:type_name -> chat.User
	0, // 1: chat.Message.type:type_name -> chat.Message.MessageType
	4, // 2: chat.Message.sent_at:type_name -> google.protobuf.Timestamp
	1, // 3: chat.MessageStream.message:type_name -> chat.Message
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_message_message_proto_init() }
func file_message_message_proto_init() {
	if File_message_message_proto != nil {
		return
	}
	file_user_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_message_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_message_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageStream); i {
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
	file_message_message_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_message_proto_goTypes,
		DependencyIndexes: file_message_message_proto_depIdxs,
		EnumInfos:         file_message_message_proto_enumTypes,
		MessageInfos:      file_message_message_proto_msgTypes,
	}.Build()
	File_message_message_proto = out.File
	file_message_message_proto_rawDesc = nil
	file_message_message_proto_goTypes = nil
	file_message_message_proto_depIdxs = nil
}
