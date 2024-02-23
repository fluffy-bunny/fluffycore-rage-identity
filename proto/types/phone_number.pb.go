// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: proto/types/phone_number.proto

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/descriptorpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PhoneType int32

const (
	PhoneType_PHONE_TYPE_UNSPECIFIED PhoneType = 0
	PhoneType_PHONE_TYPE_MOBILE      PhoneType = 1
	PhoneType_PHONE_TYPE_HOME        PhoneType = 2
	PhoneType_PHONE_TYPE_WORK        PhoneType = 3
)

// Enum value maps for PhoneType.
var (
	PhoneType_name = map[int32]string{
		0: "PHONE_TYPE_UNSPECIFIED",
		1: "PHONE_TYPE_MOBILE",
		2: "PHONE_TYPE_HOME",
		3: "PHONE_TYPE_WORK",
	}
	PhoneType_value = map[string]int32{
		"PHONE_TYPE_UNSPECIFIED": 0,
		"PHONE_TYPE_MOBILE":      1,
		"PHONE_TYPE_HOME":        2,
		"PHONE_TYPE_WORK":        3,
	}
)

func (x PhoneType) Enum() *PhoneType {
	p := new(PhoneType)
	*p = x
	return p
}

func (x PhoneType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PhoneType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_types_phone_number_proto_enumTypes[0].Descriptor()
}

func (PhoneType) Type() protoreflect.EnumType {
	return &file_proto_types_phone_number_proto_enumTypes[0]
}

func (x PhoneType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PhoneType.Descriptor instead.
func (PhoneType) EnumDescriptor() ([]byte, []int) {
	return file_proto_types_phone_number_proto_rawDescGZIP(), []int{0}
}

type PhoneTypeValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value PhoneType `protobuf:"varint,1,opt,name=value,proto3,enum=proto.types.PhoneType" json:"value,omitempty"`
}

func (x *PhoneTypeValue) Reset() {
	*x = PhoneTypeValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_phone_number_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhoneTypeValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneTypeValue) ProtoMessage() {}

func (x *PhoneTypeValue) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_phone_number_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneTypeValue.ProtoReflect.Descriptor instead.
func (*PhoneTypeValue) Descriptor() ([]byte, []int) {
	return file_proto_types_phone_number_proto_rawDescGZIP(), []int{0}
}

func (x *PhoneTypeValue) GetValue() PhoneType {
	if x != nil {
		return x.Value
	}
	return PhoneType_PHONE_TYPE_UNSPECIFIED
}

type PhoneNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CountryCode uint32    `protobuf:"varint,1,opt,name=countryCode,proto3" json:"countryCode,omitempty"`
	Number      string    `protobuf:"bytes,2,opt,name=number,proto3" json:"number,omitempty"`
	Type        PhoneType `protobuf:"varint,3,opt,name=type,proto3,enum=proto.types.PhoneType" json:"type,omitempty"`
}

func (x *PhoneNumber) Reset() {
	*x = PhoneNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_phone_number_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhoneNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneNumber) ProtoMessage() {}

func (x *PhoneNumber) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_phone_number_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneNumber.ProtoReflect.Descriptor instead.
func (*PhoneNumber) Descriptor() ([]byte, []int) {
	return file_proto_types_phone_number_proto_rawDescGZIP(), []int{1}
}

func (x *PhoneNumber) GetCountryCode() uint32 {
	if x != nil {
		return x.CountryCode
	}
	return 0
}

func (x *PhoneNumber) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *PhoneNumber) GetType() PhoneType {
	if x != nil {
		return x.Type
	}
	return PhoneType_PHONE_TYPE_UNSPECIFIED
}

type PhoneNumberDTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CountryCode uint32    `protobuf:"varint,2,opt,name=countryCode,proto3" json:"countryCode,omitempty"`
	Number      string    `protobuf:"bytes,3,opt,name=number,proto3" json:"number,omitempty"`
	Type        PhoneType `protobuf:"varint,4,opt,name=type,proto3,enum=proto.types.PhoneType" json:"type,omitempty"`
}

func (x *PhoneNumberDTO) Reset() {
	*x = PhoneNumberDTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_phone_number_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhoneNumberDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneNumberDTO) ProtoMessage() {}

func (x *PhoneNumberDTO) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_phone_number_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneNumberDTO.ProtoReflect.Descriptor instead.
func (*PhoneNumberDTO) Descriptor() ([]byte, []int) {
	return file_proto_types_phone_number_proto_rawDescGZIP(), []int{2}
}

func (x *PhoneNumberDTO) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PhoneNumberDTO) GetCountryCode() uint32 {
	if x != nil {
		return x.CountryCode
	}
	return 0
}

func (x *PhoneNumberDTO) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *PhoneNumberDTO) GetType() PhoneType {
	if x != nil {
		return x.Type
	}
	return PhoneType_PHONE_TYPE_UNSPECIFIED
}

type PhoneNumberUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CountryCode *wrapperspb.UInt32Value `protobuf:"bytes,2,opt,name=countryCode,proto3" json:"countryCode,omitempty"`
	Number      *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=number,proto3" json:"number,omitempty"`
	Type        *PhoneTypeValue         `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *PhoneNumberUpdate) Reset() {
	*x = PhoneNumberUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_phone_number_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhoneNumberUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneNumberUpdate) ProtoMessage() {}

func (x *PhoneNumberUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_phone_number_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneNumberUpdate.ProtoReflect.Descriptor instead.
func (*PhoneNumberUpdate) Descriptor() ([]byte, []int) {
	return file_proto_types_phone_number_proto_rawDescGZIP(), []int{3}
}

func (x *PhoneNumberUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PhoneNumberUpdate) GetCountryCode() *wrapperspb.UInt32Value {
	if x != nil {
		return x.CountryCode
	}
	return nil
}

func (x *PhoneNumberUpdate) GetNumber() *wrapperspb.StringValue {
	if x != nil {
		return x.Number
	}
	return nil
}

func (x *PhoneNumberUpdate) GetType() *PhoneTypeValue {
	if x != nil {
		return x.Type
	}
	return nil
}

var File_proto_types_phone_number_proto protoreflect.FileDescriptor

var file_proto_types_phone_number_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x68,
	0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x1a, 0x20, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x3e, 0x0a, 0x0e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x50,
	0x68, 0x6f, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0x73, 0x0a, 0x0b, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x20,
	0x0a, 0x0b, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0b, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x22, 0x86, 0x01, 0x0a, 0x0e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x44, 0x54, 0x4f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x72, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x2a, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x50, 0x68,
	0x6f, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xca, 0x01,
	0x0a, 0x11, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x3e, 0x0a, 0x0b, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x6f,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x33,
	0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x2a, 0x68, 0x0a, 0x09, 0x50, 0x68,
	0x6f, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x50, 0x48, 0x4f, 0x4e, 0x45,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x48, 0x4f, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x4d, 0x4f, 0x42, 0x49, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x48,
	0x4f, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x4f, 0x4d, 0x45, 0x10, 0x02, 0x12,
	0x13, 0x0a, 0x0f, 0x50, 0x48, 0x4f, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x57, 0x4f,
	0x52, 0x4b, 0x10, 0x03, 0x42, 0x7f, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x66,
	0x66, 0x79, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x72, 0x61, 0x67, 0x65, 0x6f, 0x69, 0x64, 0x63,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x50, 0x01, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x2d, 0x62, 0x75, 0x6e, 0x6e,
	0x79, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x63, 0x6f, 0x72, 0x65, 0x2d, 0x72, 0x61, 0x67,
	0x65, 0x2d, 0x6f, 0x69, 0x64, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x73, 0xaa, 0x02, 0x1a, 0x46, 0x6c, 0x75, 0x66, 0x66,
	0x79, 0x42, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x52, 0x61, 0x67, 0x65, 0x4f, 0x69, 0x64, 0x63, 0x2e,
	0x54, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_types_phone_number_proto_rawDescOnce sync.Once
	file_proto_types_phone_number_proto_rawDescData = file_proto_types_phone_number_proto_rawDesc
)

func file_proto_types_phone_number_proto_rawDescGZIP() []byte {
	file_proto_types_phone_number_proto_rawDescOnce.Do(func() {
		file_proto_types_phone_number_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_types_phone_number_proto_rawDescData)
	})
	return file_proto_types_phone_number_proto_rawDescData
}

var file_proto_types_phone_number_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_types_phone_number_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_types_phone_number_proto_goTypes = []interface{}{
	(PhoneType)(0),                 // 0: proto.types.PhoneType
	(*PhoneTypeValue)(nil),         // 1: proto.types.PhoneTypeValue
	(*PhoneNumber)(nil),            // 2: proto.types.PhoneNumber
	(*PhoneNumberDTO)(nil),         // 3: proto.types.PhoneNumberDTO
	(*PhoneNumberUpdate)(nil),      // 4: proto.types.PhoneNumberUpdate
	(*wrapperspb.UInt32Value)(nil), // 5: google.protobuf.UInt32Value
	(*wrapperspb.StringValue)(nil), // 6: google.protobuf.StringValue
}
var file_proto_types_phone_number_proto_depIdxs = []int32{
	0, // 0: proto.types.PhoneTypeValue.value:type_name -> proto.types.PhoneType
	0, // 1: proto.types.PhoneNumber.type:type_name -> proto.types.PhoneType
	0, // 2: proto.types.PhoneNumberDTO.type:type_name -> proto.types.PhoneType
	5, // 3: proto.types.PhoneNumberUpdate.countryCode:type_name -> google.protobuf.UInt32Value
	6, // 4: proto.types.PhoneNumberUpdate.number:type_name -> google.protobuf.StringValue
	1, // 5: proto.types.PhoneNumberUpdate.type:type_name -> proto.types.PhoneTypeValue
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_proto_types_phone_number_proto_init() }
func file_proto_types_phone_number_proto_init() {
	if File_proto_types_phone_number_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_types_phone_number_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhoneTypeValue); i {
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
		file_proto_types_phone_number_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhoneNumber); i {
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
		file_proto_types_phone_number_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhoneNumberDTO); i {
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
		file_proto_types_phone_number_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhoneNumberUpdate); i {
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
			RawDescriptor: file_proto_types_phone_number_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_types_phone_number_proto_goTypes,
		DependencyIndexes: file_proto_types_phone_number_proto_depIdxs,
		EnumInfos:         file_proto_types_phone_number_proto_enumTypes,
		MessageInfos:      file_proto_types_phone_number_proto_msgTypes,
	}.Build()
	File_proto_types_phone_number_proto = out.File
	file_proto_types_phone_number_proto_rawDesc = nil
	file_proto_types_phone_number_proto_goTypes = nil
	file_proto_types_phone_number_proto_depIdxs = nil
}