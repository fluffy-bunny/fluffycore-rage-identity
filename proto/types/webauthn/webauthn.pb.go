// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.1
// source: proto/types/webauthn/webauthn.proto

package webauthn

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

type Credential struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// user friendly name to identity the credential
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// A probabilistically-unique byte sequence identifying a public key credential source and its authentication assertions.
	ID []byte `protobuf:"bytes,2,opt,name=i_d,json=iD,proto3" json:"i_d,omitempty"`
	// The public key portion of a Relying Party-specific credential key pair, generated by an authenticator and returned to
	// a Relying Party at registration time (see also public key credential). The private key portion of the credential key
	// pair is known as the credential private key. Note that in the case of self attestation, the credential key pair is also
	// used as the attestation key pair, see self attestation for details.
	PublicKey []byte `protobuf:"bytes,3,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	// The attestation format used (if any) by the authenticator when creating the credential.
	AttestationType string `protobuf:"bytes,4,opt,name=attestation_type,json=attestationType,proto3" json:"attestation_type,omitempty"`
	// The transport types the authenticator supports.
	Transport []string `protobuf:"bytes,5,rep,name=transport,proto3" json:"transport,omitempty"`
	// The commonly stored flags.
	Flags *CredentialFlags `protobuf:"bytes,6,opt,name=flags,proto3" json:"flags,omitempty"`
	// The Authenticator information for a given certificate
	Authenticator *Authenticator `protobuf:"bytes,7,opt,name=authenticator,proto3" json:"authenticator,omitempty"`
}

func (x *Credential) Reset() {
	*x = Credential{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Credential) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Credential) ProtoMessage() {}

func (x *Credential) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Credential.ProtoReflect.Descriptor instead.
func (*Credential) Descriptor() ([]byte, []int) {
	return file_proto_types_webauthn_webauthn_proto_rawDescGZIP(), []int{0}
}

func (x *Credential) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Credential) GetID() []byte {
	if x != nil {
		return x.ID
	}
	return nil
}

func (x *Credential) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *Credential) GetAttestationType() string {
	if x != nil {
		return x.AttestationType
	}
	return ""
}

func (x *Credential) GetTransport() []string {
	if x != nil {
		return x.Transport
	}
	return nil
}

func (x *Credential) GetFlags() *CredentialFlags {
	if x != nil {
		return x.Flags
	}
	return nil
}

func (x *Credential) GetAuthenticator() *Authenticator {
	if x != nil {
		return x.Authenticator
	}
	return nil
}

type CredentialFlags struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Flag UP indicates the users presence.
	UserPresent bool `protobuf:"varint,1,opt,name=user_present,json=userPresent,proto3" json:"user_present,omitempty"`
	// Flag UV indicates the user verified.
	UserVerified bool `protobuf:"varint,2,opt,name=user_verified,json=userVerified,proto3" json:"user_verified,omitempty"`
	// Flag BE indicates the credential is able to be backed up and/or sync'd between devices. This should NEVER change.
	BackupEligible bool `protobuf:"varint,3,opt,name=backup_eligible,json=backupEligible,proto3" json:"backup_eligible,omitempty"`
	// Flag BS indicates the credential has been backed up and/or sync'd. This value can change but it's recommended
	// that RP's keep track of this value.
	BackupState bool `protobuf:"varint,4,opt,name=backup_state,json=backupState,proto3" json:"backup_state,omitempty"`
}

func (x *CredentialFlags) Reset() {
	*x = CredentialFlags{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CredentialFlags) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CredentialFlags) ProtoMessage() {}

func (x *CredentialFlags) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CredentialFlags.ProtoReflect.Descriptor instead.
func (*CredentialFlags) Descriptor() ([]byte, []int) {
	return file_proto_types_webauthn_webauthn_proto_rawDescGZIP(), []int{1}
}

func (x *CredentialFlags) GetUserPresent() bool {
	if x != nil {
		return x.UserPresent
	}
	return false
}

func (x *CredentialFlags) GetUserVerified() bool {
	if x != nil {
		return x.UserVerified
	}
	return false
}

func (x *CredentialFlags) GetBackupEligible() bool {
	if x != nil {
		return x.BackupEligible
	}
	return false
}

func (x *CredentialFlags) GetBackupState() bool {
	if x != nil {
		return x.BackupState
	}
	return false
}

type Authenticator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The AAGUID of the authenticator. An AAGUID is defined as an array containing the globally unique
	// identifier of the authenticator model being sought.
	AAGUID []byte `protobuf:"bytes,1,opt,name=a_a_g_u_i_d,json=aAGUID,proto3" json:"a_a_g_u_i_d,omitempty"`
	// SignCount -Upon a new login operation, the Relying Party compares the stored signature counter value
	// with the new sign_count value returned in the assertion’s authenticator data. If this new
	// signCount value is less than or equal to the stored value, a cloned authenticator may
	// exist, or the authenticator may be malfunctioning.
	SignCount uint32 `protobuf:"varint,2,opt,name=sign_count,json=signCount,proto3" json:"sign_count,omitempty"`
	// CloneWarning - This is a signal that the authenticator may be cloned, i.e. at least two copies of the
	// credential private key may exist and are being used in parallel. Relying Parties should incorporate
	// this information into their risk scoring. Whether the Relying Party updates the stored signature
	// counter value in this case, or not, or fails the authentication ceremony or not, is Relying Party-specific.
	CloneWarning bool `protobuf:"varint,3,opt,name=clone_warning,json=cloneWarning,proto3" json:"clone_warning,omitempty"`
	// Attachment is the authenticatorAttachment value returned by the request.
	Attachment   string `protobuf:"bytes,4,opt,name=attachment,proto3" json:"attachment,omitempty"`
	FriendlyName string `protobuf:"bytes,5,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty"`
}

func (x *Authenticator) Reset() {
	*x = Authenticator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Authenticator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Authenticator) ProtoMessage() {}

func (x *Authenticator) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Authenticator.ProtoReflect.Descriptor instead.
func (*Authenticator) Descriptor() ([]byte, []int) {
	return file_proto_types_webauthn_webauthn_proto_rawDescGZIP(), []int{2}
}

func (x *Authenticator) GetAAGUID() []byte {
	if x != nil {
		return x.AAGUID
	}
	return nil
}

func (x *Authenticator) GetSignCount() uint32 {
	if x != nil {
		return x.SignCount
	}
	return 0
}

func (x *Authenticator) GetCloneWarning() bool {
	if x != nil {
		return x.CloneWarning
	}
	return false
}

func (x *Authenticator) GetAttachment() string {
	if x != nil {
		return x.Attachment
	}
	return ""
}

func (x *Authenticator) GetFriendlyName() string {
	if x != nil {
		return x.FriendlyName
	}
	return ""
}

type CredentialArrayUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Update:
	//
	//	*CredentialArrayUpdate_Granular_
	//	*CredentialArrayUpdate_DeleteAll
	Update isCredentialArrayUpdate_Update `protobuf_oneof:"update"`
}

func (x *CredentialArrayUpdate) Reset() {
	*x = CredentialArrayUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CredentialArrayUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CredentialArrayUpdate) ProtoMessage() {}

func (x *CredentialArrayUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CredentialArrayUpdate.ProtoReflect.Descriptor instead.
func (*CredentialArrayUpdate) Descriptor() ([]byte, []int) {
	return file_proto_types_webauthn_webauthn_proto_rawDescGZIP(), []int{3}
}

func (m *CredentialArrayUpdate) GetUpdate() isCredentialArrayUpdate_Update {
	if m != nil {
		return m.Update
	}
	return nil
}

func (x *CredentialArrayUpdate) GetGranular() *CredentialArrayUpdate_Granular {
	if x, ok := x.GetUpdate().(*CredentialArrayUpdate_Granular_); ok {
		return x.Granular
	}
	return nil
}

func (x *CredentialArrayUpdate) GetDeleteAll() *wrapperspb.BoolValue {
	if x, ok := x.GetUpdate().(*CredentialArrayUpdate_DeleteAll); ok {
		return x.DeleteAll
	}
	return nil
}

type isCredentialArrayUpdate_Update interface {
	isCredentialArrayUpdate_Update()
}

type CredentialArrayUpdate_Granular_ struct {
	Granular *CredentialArrayUpdate_Granular `protobuf:"bytes,1,opt,name=granular,proto3,oneof"`
}

type CredentialArrayUpdate_DeleteAll struct {
	DeleteAll *wrapperspb.BoolValue `protobuf:"bytes,2,opt,name=delete_all,json=deleteAll,proto3,oneof"`
}

func (*CredentialArrayUpdate_Granular_) isCredentialArrayUpdate_Update() {}

func (*CredentialArrayUpdate_DeleteAll) isCredentialArrayUpdate_Update() {}

type CredentialArrayUpdate_Granular struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Add           []*Credential `protobuf:"bytes,1,rep,name=add,proto3" json:"add,omitempty"`
	RemoveAAGUIDs [][]byte      `protobuf:"bytes,2,rep,name=remove_a_a_g_u_i_ds,json=removeAAGUIDs,proto3" json:"remove_a_a_g_u_i_ds,omitempty"`
}

func (x *CredentialArrayUpdate_Granular) Reset() {
	*x = CredentialArrayUpdate_Granular{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CredentialArrayUpdate_Granular) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CredentialArrayUpdate_Granular) ProtoMessage() {}

func (x *CredentialArrayUpdate_Granular) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_webauthn_webauthn_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CredentialArrayUpdate_Granular.ProtoReflect.Descriptor instead.
func (*CredentialArrayUpdate_Granular) Descriptor() ([]byte, []int) {
	return file_proto_types_webauthn_webauthn_proto_rawDescGZIP(), []int{3, 0}
}

func (x *CredentialArrayUpdate_Granular) GetAdd() []*Credential {
	if x != nil {
		return x.Add
	}
	return nil
}

func (x *CredentialArrayUpdate_Granular) GetRemoveAAGUIDs() [][]byte {
	if x != nil {
		return x.RemoveAAGUIDs
	}
	return nil
}

var File_proto_types_webauthn_webauthn_proto protoreflect.FileDescriptor

var file_proto_types_webauthn_webauthn_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x77, 0x65,
	0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2f, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x1a, 0x20, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x02,
	0x0a, 0x0a, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x0f, 0x0a, 0x03, 0x69, 0x5f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69,
	0x44, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79,
	0x12, 0x29, 0x0a, 0x10, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x74, 0x74, 0x65,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x3b, 0x0a, 0x05, 0x66, 0x6c, 0x61,
	0x67, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2e,
	0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x52,
	0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x49, 0x0a, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e,
	0x74, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x77, 0x65, 0x62, 0x61,
	0x75, 0x74, 0x68, 0x6e, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x52, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f,
	0x72, 0x22, 0xa5, 0x01, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70, 0x72,
	0x65, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x75, 0x73, 0x65,
	0x72, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0c, 0x75, 0x73, 0x65, 0x72, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x27, 0x0a,
	0x0f, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x5f, 0x65, 0x6c, 0x69, 0x67, 0x69, 0x62, 0x6c, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x45, 0x6c,
	0x69, 0x67, 0x69, 0x62, 0x6c, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x62, 0x61,
	0x63, 0x6b, 0x75, 0x70, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0xb5, 0x01, 0x0a, 0x0d, 0x41, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x0b, 0x61,
	0x5f, 0x61, 0x5f, 0x67, 0x5f, 0x75, 0x5f, 0x69, 0x5f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x06, 0x61, 0x41, 0x47, 0x55, 0x49, 0x44, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x73, 0x69,
	0x67, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x6f, 0x6e, 0x65,
	0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x63, 0x6c, 0x6f, 0x6e, 0x65, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x12, 0x1e, 0x0a, 0x0a,
	0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d,
	0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0x9e, 0x02, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x41, 0x72, 0x72, 0x61, 0x79, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x52, 0x0a, 0x08, 0x67,
	0x72, 0x61, 0x6e, 0x75, 0x6c, 0x61, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x77, 0x65, 0x62, 0x61,
	0x75, 0x74, 0x68, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x41,
	0x72, 0x72, 0x61, 0x79, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x47, 0x72, 0x61, 0x6e, 0x75,
	0x6c, 0x61, 0x72, 0x48, 0x00, 0x52, 0x08, 0x67, 0x72, 0x61, 0x6e, 0x75, 0x6c, 0x61, 0x72, 0x12,
	0x3b, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x61, 0x6c, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48,
	0x00, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x1a, 0x6a, 0x0a, 0x08,
	0x47, 0x72, 0x61, 0x6e, 0x75, 0x6c, 0x61, 0x72, 0x12, 0x32, 0x0a, 0x03, 0x61, 0x64, 0x64, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x2e, 0x43, 0x72, 0x65,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x52, 0x03, 0x61, 0x64, 0x64, 0x12, 0x2a, 0x0a, 0x13,
	0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x5f, 0x61, 0x5f, 0x61, 0x5f, 0x67, 0x5f, 0x75, 0x5f, 0x69,
	0x5f, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0d, 0x72, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x41, 0x41, 0x47, 0x55, 0x49, 0x44, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x42, 0x50, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x2d, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2f, 0x66, 0x6c,
	0x75, 0x66, 0x66, 0x79, 0x63, 0x6f, 0x72, 0x65, 0x2d, 0x72, 0x61, 0x67, 0x65, 0x2d, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2f, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x3b, 0x77, 0x65, 0x62, 0x61,
	0x75, 0x74, 0x68, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_types_webauthn_webauthn_proto_rawDescOnce sync.Once
	file_proto_types_webauthn_webauthn_proto_rawDescData = file_proto_types_webauthn_webauthn_proto_rawDesc
)

func file_proto_types_webauthn_webauthn_proto_rawDescGZIP() []byte {
	file_proto_types_webauthn_webauthn_proto_rawDescOnce.Do(func() {
		file_proto_types_webauthn_webauthn_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_types_webauthn_webauthn_proto_rawDescData)
	})
	return file_proto_types_webauthn_webauthn_proto_rawDescData
}

var file_proto_types_webauthn_webauthn_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_types_webauthn_webauthn_proto_goTypes = []interface{}{
	(*Credential)(nil),                     // 0: proto.types.webauthn.Credential
	(*CredentialFlags)(nil),                // 1: proto.types.webauthn.CredentialFlags
	(*Authenticator)(nil),                  // 2: proto.types.webauthn.Authenticator
	(*CredentialArrayUpdate)(nil),          // 3: proto.types.webauthn.CredentialArrayUpdate
	(*CredentialArrayUpdate_Granular)(nil), // 4: proto.types.webauthn.CredentialArrayUpdate.Granular
	(*wrapperspb.BoolValue)(nil),           // 5: google.protobuf.BoolValue
}
var file_proto_types_webauthn_webauthn_proto_depIdxs = []int32{
	1, // 0: proto.types.webauthn.Credential.flags:type_name -> proto.types.webauthn.CredentialFlags
	2, // 1: proto.types.webauthn.Credential.authenticator:type_name -> proto.types.webauthn.Authenticator
	4, // 2: proto.types.webauthn.CredentialArrayUpdate.granular:type_name -> proto.types.webauthn.CredentialArrayUpdate.Granular
	5, // 3: proto.types.webauthn.CredentialArrayUpdate.delete_all:type_name -> google.protobuf.BoolValue
	0, // 4: proto.types.webauthn.CredentialArrayUpdate.Granular.add:type_name -> proto.types.webauthn.Credential
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_types_webauthn_webauthn_proto_init() }
func file_proto_types_webauthn_webauthn_proto_init() {
	if File_proto_types_webauthn_webauthn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_types_webauthn_webauthn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Credential); i {
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
		file_proto_types_webauthn_webauthn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CredentialFlags); i {
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
		file_proto_types_webauthn_webauthn_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Authenticator); i {
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
		file_proto_types_webauthn_webauthn_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CredentialArrayUpdate); i {
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
		file_proto_types_webauthn_webauthn_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CredentialArrayUpdate_Granular); i {
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
	file_proto_types_webauthn_webauthn_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*CredentialArrayUpdate_Granular_)(nil),
		(*CredentialArrayUpdate_DeleteAll)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_types_webauthn_webauthn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_types_webauthn_webauthn_proto_goTypes,
		DependencyIndexes: file_proto_types_webauthn_webauthn_proto_depIdxs,
		MessageInfos:      file_proto_types_webauthn_webauthn_proto_msgTypes,
	}.Build()
	File_proto_types_webauthn_webauthn_proto = out.File
	file_proto_types_webauthn_webauthn_proto_rawDesc = nil
	file_proto_types_webauthn_webauthn_proto_goTypes = nil
	file_proto_types_webauthn_webauthn_proto_depIdxs = nil
}
