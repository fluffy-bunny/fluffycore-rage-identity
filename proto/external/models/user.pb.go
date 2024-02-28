// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.25.1
// source: proto/external/models/user.proto

package models

import (
	types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
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

type UserState int32

const (
	UserState_USER_STATE_UNSPECIFIED UserState = 0
	UserState_USER_STATE_ACTIVE      UserState = 1
	UserState_USER_STATE_DISABLED    UserState = 2
	UserState_USER_STATE_DELETED     UserState = 3
	UserState_USER_STATE_PENDING     UserState = 4
)

// Enum value maps for UserState.
var (
	UserState_name = map[int32]string{
		0: "USER_STATE_UNSPECIFIED",
		1: "USER_STATE_ACTIVE",
		2: "USER_STATE_DISABLED",
		3: "USER_STATE_DELETED",
		4: "USER_STATE_PENDING",
	}
	UserState_value = map[string]int32{
		"USER_STATE_UNSPECIFIED": 0,
		"USER_STATE_ACTIVE":      1,
		"USER_STATE_DISABLED":    2,
		"USER_STATE_DELETED":     3,
		"USER_STATE_PENDING":     4,
	}
)

func (x UserState) Enum() *UserState {
	p := new(UserState)
	*p = x
	return p
}

func (x UserState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserState) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_external_models_user_proto_enumTypes[0].Descriptor()
}

func (UserState) Type() protoreflect.EnumType {
	return &file_proto_external_models_user_proto_enumTypes[0]
}

func (x UserState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserState.Descriptor instead.
func (UserState) EnumDescriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{0}
}

type UserStateFilterExpression struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eq UserState   `protobuf:"varint,1,opt,name=eq,proto3,enum=proto.external.models.UserState" json:"eq,omitempty"`
	In []UserState `protobuf:"varint,2,rep,packed,name=in,proto3,enum=proto.external.models.UserState" json:"in,omitempty"`
}

func (x *UserStateFilterExpression) Reset() {
	*x = UserStateFilterExpression{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserStateFilterExpression) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserStateFilterExpression) ProtoMessage() {}

func (x *UserStateFilterExpression) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserStateFilterExpression.ProtoReflect.Descriptor instead.
func (*UserStateFilterExpression) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{0}
}

func (x *UserStateFilterExpression) GetEq() UserState {
	if x != nil {
		return x.Eq
	}
	return UserState_USER_STATE_UNSPECIFIED
}

func (x *UserStateFilterExpression) GetIn() []UserState {
	if x != nil {
		return x.In
	}
	return nil
}

type UserStateValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value UserState `protobuf:"varint,1,opt,name=value,proto3,enum=proto.external.models.UserState" json:"value,omitempty"`
}

func (x *UserStateValue) Reset() {
	*x = UserStateValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserStateValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserStateValue) ProtoMessage() {}

func (x *UserStateValue) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserStateValue.ProtoReflect.Descriptor instead.
func (*UserStateValue) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{1}
}

func (x *UserStateValue) GetValue() UserState {
	if x != nil {
		return x.Value
	}
	return UserState_USER_STATE_UNSPECIFIED
}

type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Street     string `protobuf:"bytes,1,opt,name=street,proto3" json:"street,omitempty"`
	City       string `protobuf:"bytes,2,opt,name=city,proto3" json:"city,omitempty"`
	State      string `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
	PostalCode string `protobuf:"bytes,4,opt,name=postal_code,json=postalCode,proto3" json:"postal_code,omitempty"`
	Country    string `protobuf:"bytes,5,opt,name=country,proto3" json:"country,omitempty"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{2}
}

func (x *Address) GetStreet() string {
	if x != nil {
		return x.Street
	}
	return ""
}

func (x *Address) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *Address) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *Address) GetPostalCode() string {
	if x != nil {
		return x.PostalCode
	}
	return ""
}

func (x *Address) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

type AddressUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Street     *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=street,proto3" json:"street,omitempty"`
	City       *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=city,proto3" json:"city,omitempty"`
	State      *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
	PostalCode *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=postal_code,json=postalCode,proto3" json:"postal_code,omitempty"`
	Country    *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=country,proto3" json:"country,omitempty"`
}

func (x *AddressUpdate) Reset() {
	*x = AddressUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressUpdate) ProtoMessage() {}

func (x *AddressUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressUpdate.ProtoReflect.Descriptor instead.
func (*AddressUpdate) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{3}
}

func (x *AddressUpdate) GetStreet() *wrapperspb.StringValue {
	if x != nil {
		return x.Street
	}
	return nil
}

func (x *AddressUpdate) GetCity() *wrapperspb.StringValue {
	if x != nil {
		return x.City
	}
	return nil
}

func (x *AddressUpdate) GetState() *wrapperspb.StringValue {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *AddressUpdate) GetPostalCode() *wrapperspb.StringValue {
	if x != nil {
		return x.PostalCode
	}
	return nil
}

func (x *AddressUpdate) GetCountry() *wrapperspb.StringValue {
	if x != nil {
		return x.Country
	}
	return nil
}

type Profile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GivenName    string                  `protobuf:"bytes,1,opt,name=given_name,json=givenName,proto3" json:"given_name,omitempty"`
	FamilyName   string                  `protobuf:"bytes,2,opt,name=family_name,json=familyName,proto3" json:"family_name,omitempty"`
	PhoneNumbers []*types.PhoneNumberDTO `protobuf:"bytes,3,rep,name=phone_numbers,json=phoneNumbers,proto3" json:"phone_numbers,omitempty"`
	Address      *Address                `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Profile) Reset() {
	*x = Profile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profile) ProtoMessage() {}

func (x *Profile) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profile.ProtoReflect.Descriptor instead.
func (*Profile) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{4}
}

func (x *Profile) GetGivenName() string {
	if x != nil {
		return x.GivenName
	}
	return ""
}

func (x *Profile) GetFamilyName() string {
	if x != nil {
		return x.FamilyName
	}
	return ""
}

func (x *Profile) GetPhoneNumbers() []*types.PhoneNumberDTO {
	if x != nil {
		return x.PhoneNumbers
	}
	return nil
}

func (x *Profile) GetAddress() *Address {
	if x != nil {
		return x.Address
	}
	return nil
}

type ProfileUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GivenName    *wrapperspb.StringValue       `protobuf:"bytes,1,opt,name=given_name,json=givenName,proto3" json:"given_name,omitempty"`
	FamilyName   *wrapperspb.StringValue       `protobuf:"bytes,2,opt,name=family_name,json=familyName,proto3" json:"family_name,omitempty"`
	PhoneNumbers []*types.PhoneNumberDTOUpdate `protobuf:"bytes,3,rep,name=phone_numbers,json=phoneNumbers,proto3" json:"phone_numbers,omitempty"`
	Address      *AddressUpdate                `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *ProfileUpdate) Reset() {
	*x = ProfileUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileUpdate) ProtoMessage() {}

func (x *ProfileUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileUpdate.ProtoReflect.Descriptor instead.
func (*ProfileUpdate) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{5}
}

func (x *ProfileUpdate) GetGivenName() *wrapperspb.StringValue {
	if x != nil {
		return x.GivenName
	}
	return nil
}

func (x *ProfileUpdate) GetFamilyName() *wrapperspb.StringValue {
	if x != nil {
		return x.FamilyName
	}
	return nil
}

func (x *ProfileUpdate) GetPhoneNumbers() []*types.PhoneNumberDTOUpdate {
	if x != nil {
		return x.PhoneNumbers
	}
	return nil
}

func (x *ProfileUpdate) GetAddress() *AddressUpdate {
	if x != nil {
		return x.Address
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Metadata map[string]string `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	State    UserState         `protobuf:"varint,3,opt,name=state,proto3,enum=proto.external.models.UserState" json:"state,omitempty"`
	Profile  *Profile          `protobuf:"bytes,4,opt,name=profile,proto3" json:"profile,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{6}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *User) GetState() UserState {
	if x != nil {
		return x.State
	}
	return UserState_USER_STATE_UNSPECIFIED
}

func (x *User) GetProfile() *Profile {
	if x != nil {
		return x.Profile
	}
	return nil
}

type UserUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Metadata *types.StringMapUpdate `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	State    *UserStateValue        `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
	Profile  *ProfileUpdate         `protobuf:"bytes,4,opt,name=profile,proto3" json:"profile,omitempty"`
}

func (x *UserUpdate) Reset() {
	*x = UserUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_external_models_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUpdate) ProtoMessage() {}

func (x *UserUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_external_models_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUpdate.ProtoReflect.Descriptor instead.
func (*UserUpdate) Descriptor() ([]byte, []int) {
	return file_proto_external_models_user_proto_rawDescGZIP(), []int{7}
}

func (x *UserUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UserUpdate) GetMetadata() *types.StringMapUpdate {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *UserUpdate) GetState() *UserStateValue {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *UserUpdate) GetProfile() *ProfileUpdate {
	if x != nil {
		return x.Profile
	}
	return nil
}

var File_proto_external_models_user_proto protoreflect.FileDescriptor

var file_proto_external_models_user_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69,
	0x76, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x19, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x36, 0x0a, 0x02, 0x65, 0x71, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x42,
	0x04, 0xc8, 0xc0, 0x29, 0x0a, 0x52, 0x02, 0x65, 0x71, 0x12, 0x36, 0x0a, 0x02, 0x69, 0x6e, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x42, 0x04, 0xc8, 0xc0, 0x29, 0x21, 0x52, 0x02, 0x69,
	0x6e, 0x3a, 0x04, 0xe0, 0xc6, 0x29, 0x01, 0x22, 0x48, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x86, 0x01, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x74, 0x72, 0x65, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x22, 0xa2, 0x02, 0x0a, 0x0d, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x34, 0x0a, 0x06,
	0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65,
	0x65, 0x74, 0x12, 0x30, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04,
	0x63, 0x69, 0x74, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x6f, 0x73, 0x74,
	0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x70, 0x6f, 0x73,
	0x74, 0x61, 0x6c, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x22,
	0xc5, 0x01, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x67,
	0x69, 0x76, 0x65, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x67, 0x69, 0x76, 0x65, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x61,
	0x6d, 0x69, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x0d, 0x70,
	0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x54, 0x4f, 0x52,
	0x0c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x38, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x93, 0x02, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x67, 0x69, 0x76,
	0x65, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x67, 0x69, 0x76,
	0x65, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x0b, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x66, 0x61, 0x6d, 0x69, 0x6c,
	0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x46, 0x0a, 0x0d, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x54, 0x4f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x0c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x3e, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x8c, 0x02,
	0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x45, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x36, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x1a,
	0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xd3, 0x01, 0x0a,
	0x0a, 0x55, 0x73, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x38, 0x0a, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x4d, 0x61, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3b, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x2a, 0x87, 0x01, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x1a, 0x0a, 0x16, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56,
	0x45, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x45, 0x5f, 0x44, 0x49, 0x53, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54,
	0x45, 0x44, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x42, 0x4f, 0x5a, 0x4d,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66,
	0x79, 0x2d, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x63, 0x6f,
	0x72, 0x65, 0x2d, 0x72, 0x61, 0x67, 0x65, 0x2d, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_external_models_user_proto_rawDescOnce sync.Once
	file_proto_external_models_user_proto_rawDescData = file_proto_external_models_user_proto_rawDesc
)

func file_proto_external_models_user_proto_rawDescGZIP() []byte {
	file_proto_external_models_user_proto_rawDescOnce.Do(func() {
		file_proto_external_models_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_external_models_user_proto_rawDescData)
	})
	return file_proto_external_models_user_proto_rawDescData
}

var file_proto_external_models_user_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_external_models_user_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_external_models_user_proto_goTypes = []interface{}{
	(UserState)(0),                     // 0: proto.external.models.UserState
	(*UserStateFilterExpression)(nil),  // 1: proto.external.models.UserStateFilterExpression
	(*UserStateValue)(nil),             // 2: proto.external.models.UserStateValue
	(*Address)(nil),                    // 3: proto.external.models.Address
	(*AddressUpdate)(nil),              // 4: proto.external.models.AddressUpdate
	(*Profile)(nil),                    // 5: proto.external.models.Profile
	(*ProfileUpdate)(nil),              // 6: proto.external.models.ProfileUpdate
	(*User)(nil),                       // 7: proto.external.models.User
	(*UserUpdate)(nil),                 // 8: proto.external.models.UserUpdate
	nil,                                // 9: proto.external.models.User.MetadataEntry
	(*wrapperspb.StringValue)(nil),     // 10: google.protobuf.StringValue
	(*types.PhoneNumberDTO)(nil),       // 11: proto.types.PhoneNumberDTO
	(*types.PhoneNumberDTOUpdate)(nil), // 12: proto.types.PhoneNumberDTOUpdate
	(*types.StringMapUpdate)(nil),      // 13: proto.types.StringMapUpdate
}
var file_proto_external_models_user_proto_depIdxs = []int32{
	0,  // 0: proto.external.models.UserStateFilterExpression.eq:type_name -> proto.external.models.UserState
	0,  // 1: proto.external.models.UserStateFilterExpression.in:type_name -> proto.external.models.UserState
	0,  // 2: proto.external.models.UserStateValue.value:type_name -> proto.external.models.UserState
	10, // 3: proto.external.models.AddressUpdate.street:type_name -> google.protobuf.StringValue
	10, // 4: proto.external.models.AddressUpdate.city:type_name -> google.protobuf.StringValue
	10, // 5: proto.external.models.AddressUpdate.state:type_name -> google.protobuf.StringValue
	10, // 6: proto.external.models.AddressUpdate.postal_code:type_name -> google.protobuf.StringValue
	10, // 7: proto.external.models.AddressUpdate.country:type_name -> google.protobuf.StringValue
	11, // 8: proto.external.models.Profile.phone_numbers:type_name -> proto.types.PhoneNumberDTO
	3,  // 9: proto.external.models.Profile.address:type_name -> proto.external.models.Address
	10, // 10: proto.external.models.ProfileUpdate.given_name:type_name -> google.protobuf.StringValue
	10, // 11: proto.external.models.ProfileUpdate.family_name:type_name -> google.protobuf.StringValue
	12, // 12: proto.external.models.ProfileUpdate.phone_numbers:type_name -> proto.types.PhoneNumberDTOUpdate
	4,  // 13: proto.external.models.ProfileUpdate.address:type_name -> proto.external.models.AddressUpdate
	9,  // 14: proto.external.models.User.metadata:type_name -> proto.external.models.User.MetadataEntry
	0,  // 15: proto.external.models.User.state:type_name -> proto.external.models.UserState
	5,  // 16: proto.external.models.User.profile:type_name -> proto.external.models.Profile
	13, // 17: proto.external.models.UserUpdate.metadata:type_name -> proto.types.StringMapUpdate
	2,  // 18: proto.external.models.UserUpdate.state:type_name -> proto.external.models.UserStateValue
	6,  // 19: proto.external.models.UserUpdate.profile:type_name -> proto.external.models.ProfileUpdate
	20, // [20:20] is the sub-list for method output_type
	20, // [20:20] is the sub-list for method input_type
	20, // [20:20] is the sub-list for extension type_name
	20, // [20:20] is the sub-list for extension extendee
	0,  // [0:20] is the sub-list for field type_name
}

func init() { file_proto_external_models_user_proto_init() }
func file_proto_external_models_user_proto_init() {
	if File_proto_external_models_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_external_models_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserStateFilterExpression); i {
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
		file_proto_external_models_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserStateValue); i {
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
		file_proto_external_models_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
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
		file_proto_external_models_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddressUpdate); i {
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
		file_proto_external_models_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profile); i {
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
		file_proto_external_models_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileUpdate); i {
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
		file_proto_external_models_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_proto_external_models_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserUpdate); i {
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
			RawDescriptor: file_proto_external_models_user_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_external_models_user_proto_goTypes,
		DependencyIndexes: file_proto_external_models_user_proto_depIdxs,
		EnumInfos:         file_proto_external_models_user_proto_enumTypes,
		MessageInfos:      file_proto_external_models_user_proto_msgTypes,
	}.Build()
	File_proto_external_models_user_proto = out.File
	file_proto_external_models_user_proto_rawDesc = nil
	file_proto_external_models_user_proto_goTypes = nil
	file_proto_external_models_user_proto_depIdxs = nil
}
