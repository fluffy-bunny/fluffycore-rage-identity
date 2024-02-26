// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: proto/oidc/models/idp.proto

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

type OIDCProtocol struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authority    string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	ClientId     string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,3,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
	Scope        string `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
}

func (x *OIDCProtocol) Reset() {
	*x = OIDCProtocol{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OIDCProtocol) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OIDCProtocol) ProtoMessage() {}

func (x *OIDCProtocol) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OIDCProtocol.ProtoReflect.Descriptor instead.
func (*OIDCProtocol) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{0}
}

func (x *OIDCProtocol) GetAuthority() string {
	if x != nil {
		return x.Authority
	}
	return ""
}

func (x *OIDCProtocol) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *OIDCProtocol) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *OIDCProtocol) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

type GithubOAuth2Protocol struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,2,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
}

func (x *GithubOAuth2Protocol) Reset() {
	*x = GithubOAuth2Protocol{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GithubOAuth2Protocol) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GithubOAuth2Protocol) ProtoMessage() {}

func (x *GithubOAuth2Protocol) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GithubOAuth2Protocol.ProtoReflect.Descriptor instead.
func (*GithubOAuth2Protocol) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{1}
}

func (x *GithubOAuth2Protocol) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *GithubOAuth2Protocol) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

type OAuth2Protocol struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId              string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret          string `protobuf:"bytes,2,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
	Scope                 string `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	AuthorizationEndpoint string `protobuf:"bytes,4,opt,name=authorization_endpoint,json=authorizationEndpoint,proto3" json:"authorization_endpoint,omitempty"`
	TokenEndpoint         string `protobuf:"bytes,5,opt,name=token_endpoint,json=tokenEndpoint,proto3" json:"token_endpoint,omitempty"`
}

func (x *OAuth2Protocol) Reset() {
	*x = OAuth2Protocol{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OAuth2Protocol) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OAuth2Protocol) ProtoMessage() {}

func (x *OAuth2Protocol) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OAuth2Protocol.ProtoReflect.Descriptor instead.
func (*OAuth2Protocol) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{2}
}

func (x *OAuth2Protocol) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *OAuth2Protocol) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *OAuth2Protocol) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *OAuth2Protocol) GetAuthorizationEndpoint() string {
	if x != nil {
		return x.AuthorizationEndpoint
	}
	return ""
}

func (x *OAuth2Protocol) GetTokenEndpoint() string {
	if x != nil {
		return x.TokenEndpoint
	}
	return ""
}

type Protocol struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*Protocol_Oidc
	//	*Protocol_Oauth2
	//	*Protocol_Github
	Value isProtocol_Value `protobuf_oneof:"value"`
}

func (x *Protocol) Reset() {
	*x = Protocol{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Protocol) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Protocol) ProtoMessage() {}

func (x *Protocol) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Protocol.ProtoReflect.Descriptor instead.
func (*Protocol) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{3}
}

func (m *Protocol) GetValue() isProtocol_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Protocol) GetOidc() *OIDCProtocol {
	if x, ok := x.GetValue().(*Protocol_Oidc); ok {
		return x.Oidc
	}
	return nil
}

func (x *Protocol) GetOauth2() *OAuth2Protocol {
	if x, ok := x.GetValue().(*Protocol_Oauth2); ok {
		return x.Oauth2
	}
	return nil
}

func (x *Protocol) GetGithub() *GithubOAuth2Protocol {
	if x, ok := x.GetValue().(*Protocol_Github); ok {
		return x.Github
	}
	return nil
}

type isProtocol_Value interface {
	isProtocol_Value()
}

type Protocol_Oidc struct {
	Oidc *OIDCProtocol `protobuf:"bytes,1,opt,name=oidc,proto3,oneof"`
}

type Protocol_Oauth2 struct {
	Oauth2 *OAuth2Protocol `protobuf:"bytes,2,opt,name=oauth2,proto3,oneof"`
}

type Protocol_Github struct {
	Github *GithubOAuth2Protocol `protobuf:"bytes,3,opt,name=github,proto3,oneof"`
}

func (*Protocol_Oidc) isProtocol_Value() {}

func (*Protocol_Oauth2) isProtocol_Value() {}

func (*Protocol_Github) isProtocol_Value() {}

type ProtocolUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value *Protocol `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ProtocolUpdate) Reset() {
	*x = ProtocolUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtocolUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtocolUpdate) ProtoMessage() {}

func (x *ProtocolUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtocolUpdate.ProtoReflect.Descriptor instead.
func (*ProtocolUpdate) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{4}
}

func (x *ProtocolUpdate) GetValue() *Protocol {
	if x != nil {
		return x.Value
	}
	return nil
}

type IDP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                        string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Enabled                   bool              `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Slug                      string            `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Name                      string            `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Description               string            `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Protocol                  *Protocol         `protobuf:"bytes,6,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Metadata                  map[string]string `protobuf:"bytes,7,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ClaimedDomains            []string          `protobuf:"bytes,8,rep,name=claimed_domains,json=claimedDomains,proto3" json:"claimed_domains,omitempty"`
	Hidden                    bool              `protobuf:"varint,9,opt,name=hidden,proto3" json:"hidden,omitempty"`
	EmailVerificationRequired bool              `protobuf:"varint,10,opt,name=email_verification_required,json=emailVerificationRequired,proto3" json:"email_verification_required,omitempty"`
	AutoCreate                bool              `protobuf:"varint,11,opt,name=auto_create,json=autoCreate,proto3" json:"auto_create,omitempty"`
}

func (x *IDP) Reset() {
	*x = IDP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IDP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IDP) ProtoMessage() {}

func (x *IDP) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IDP.ProtoReflect.Descriptor instead.
func (*IDP) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{5}
}

func (x *IDP) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *IDP) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *IDP) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

func (x *IDP) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *IDP) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *IDP) GetProtocol() *Protocol {
	if x != nil {
		return x.Protocol
	}
	return nil
}

func (x *IDP) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *IDP) GetClaimedDomains() []string {
	if x != nil {
		return x.ClaimedDomains
	}
	return nil
}

func (x *IDP) GetHidden() bool {
	if x != nil {
		return x.Hidden
	}
	return false
}

func (x *IDP) GetEmailVerificationRequired() bool {
	if x != nil {
		return x.EmailVerificationRequired
	}
	return false
}

func (x *IDP) GetAutoCreate() bool {
	if x != nil {
		return x.AutoCreate
	}
	return false
}

type IDPs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Idps []*IDP `protobuf:"bytes,1,rep,name=idps,proto3" json:"idps,omitempty"`
}

func (x *IDPs) Reset() {
	*x = IDPs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IDPs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IDPs) ProtoMessage() {}

func (x *IDPs) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IDPs.ProtoReflect.Descriptor instead.
func (*IDPs) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{6}
}

func (x *IDPs) GetIdps() []*IDP {
	if x != nil {
		return x.Idps
	}
	return nil
}

type IDPUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                        string                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Enabled                   *wrapperspb.BoolValue    `protobuf:"bytes,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Slug                      *wrapperspb.StringValue  `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Name                      *wrapperspb.StringValue  `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Description               *wrapperspb.StringValue  `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Protocol                  *ProtocolUpdate          `protobuf:"bytes,6,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Metadata                  *types.StringMapUpdate   `protobuf:"bytes,7,opt,name=metadata,proto3" json:"metadata,omitempty"`
	ClaimedDomains            *types.StringArrayUpdate `protobuf:"bytes,8,opt,name=claimed_domains,json=claimedDomains,proto3" json:"claimed_domains,omitempty"`
	Hidden                    *wrapperspb.BoolValue    `protobuf:"bytes,9,opt,name=hidden,proto3" json:"hidden,omitempty"`
	EmailVerificationRequired *wrapperspb.BoolValue    `protobuf:"bytes,10,opt,name=email_verification_required,json=emailVerificationRequired,proto3" json:"email_verification_required,omitempty"`
	AutoCreate                *wrapperspb.BoolValue    `protobuf:"bytes,11,opt,name=auto_create,json=autoCreate,proto3" json:"auto_create,omitempty"`
}

func (x *IDPUpdate) Reset() {
	*x = IDPUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_models_idp_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IDPUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IDPUpdate) ProtoMessage() {}

func (x *IDPUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_models_idp_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IDPUpdate.ProtoReflect.Descriptor instead.
func (*IDPUpdate) Descriptor() ([]byte, []int) {
	return file_proto_oidc_models_idp_proto_rawDescGZIP(), []int{7}
}

func (x *IDPUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *IDPUpdate) GetEnabled() *wrapperspb.BoolValue {
	if x != nil {
		return x.Enabled
	}
	return nil
}

func (x *IDPUpdate) GetSlug() *wrapperspb.StringValue {
	if x != nil {
		return x.Slug
	}
	return nil
}

func (x *IDPUpdate) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *IDPUpdate) GetDescription() *wrapperspb.StringValue {
	if x != nil {
		return x.Description
	}
	return nil
}

func (x *IDPUpdate) GetProtocol() *ProtocolUpdate {
	if x != nil {
		return x.Protocol
	}
	return nil
}

func (x *IDPUpdate) GetMetadata() *types.StringMapUpdate {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *IDPUpdate) GetClaimedDomains() *types.StringArrayUpdate {
	if x != nil {
		return x.ClaimedDomains
	}
	return nil
}

func (x *IDPUpdate) GetHidden() *wrapperspb.BoolValue {
	if x != nil {
		return x.Hidden
	}
	return nil
}

func (x *IDPUpdate) GetEmailVerificationRequired() *wrapperspb.BoolValue {
	if x != nil {
		return x.EmailVerificationRequired
	}
	return nil
}

func (x *IDPUpdate) GetAutoCreate() *wrapperspb.BoolValue {
	if x != nil {
		return x.AutoCreate
	}
	return nil
}

var File_proto_oidc_models_idp_proto protoreflect.FileDescriptor

var file_proto_oidc_models_idp_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x69, 0x64, 0x63, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2f, 0x69, 0x64, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f,
	0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x84, 0x01, 0x0a, 0x0c, 0x4f, 0x49, 0x44, 0x43, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x22, 0x58, 0x0a, 0x14, 0x47, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12,
	0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x22, 0xc6, 0x01, 0x0a, 0x0e, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x35, 0x0a, 0x16,
	0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x61, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0xca, 0x01, 0x0a, 0x08, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x35, 0x0a, 0x04, 0x6f, 0x69, 0x64, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69,
	0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4f, 0x49, 0x44, 0x43, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x00, 0x52, 0x04, 0x6f, 0x69, 0x64, 0x63, 0x12, 0x3b,
	0x0a, 0x06, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x48, 0x00, 0x52, 0x06, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x12, 0x41, 0x0a, 0x06, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x00, 0x52, 0x06, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x42, 0x07,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x43, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x31, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xd3, 0x03, 0x0a,
	0x03, 0x49, 0x44, 0x50, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c,
	0x75, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x12, 0x40, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x49, 0x44, 0x50, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x65, 0x64, 0x5f, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6c,
	0x61, 0x69, 0x6d, 0x65, 0x64, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x68, 0x69,
	0x64, 0x64, 0x65, 0x6e, 0x12, 0x3e, 0x0a, 0x1b, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x19, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x69, 0x72, 0x65, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x61, 0x75, 0x74, 0x6f, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x32, 0x0a, 0x04, 0x49, 0x44, 0x50, 0x73, 0x12, 0x2a, 0x0a, 0x04, 0x69, 0x64,
	0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x49, 0x44, 0x50,
	0x52, 0x04, 0x69, 0x64, 0x70, 0x73, 0x22, 0x84, 0x05, 0x0a, 0x09, 0x49, 0x44, 0x50, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x73, 0x6c,
	0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x12, 0x30, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3e,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3d,
	0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x38, 0x0a,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x47, 0x0a, 0x0f, 0x63, 0x6c, 0x61, 0x69, 0x6d,
	0x65, 0x64, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x41, 0x72, 0x72, 0x61, 0x79, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x52, 0x0e, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x65, 0x64, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73,
	0x12, 0x32, 0x0a, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x68, 0x69,
	0x64, 0x64, 0x65, 0x6e, 0x12, 0x5a, 0x0a, 0x1b, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x19, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x12, 0x3b, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x0a, 0x61, 0x75, 0x74, 0x6f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x8d, 0x01,
	0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x62, 0x75, 0x6e, 0x6e,
	0x79, 0x2e, 0x72, 0x61, 0x67, 0x65, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x50, 0x01, 0x5a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x2d, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2f, 0x66, 0x6c, 0x75,
	0x66, 0x66, 0x79, 0x63, 0x6f, 0x72, 0x65, 0x2d, 0x72, 0x61, 0x67, 0x65, 0x2d, 0x6f, 0x69, 0x64,
	0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x69, 0x64, 0x63, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0xaa, 0x02, 0x20, 0x46, 0x6c, 0x75,
	0x66, 0x66, 0x79, 0x42, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x52, 0x61, 0x67, 0x65, 0x4f, 0x69, 0x64,
	0x63, 0x2e, 0x4f, 0x69, 0x64, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_oidc_models_idp_proto_rawDescOnce sync.Once
	file_proto_oidc_models_idp_proto_rawDescData = file_proto_oidc_models_idp_proto_rawDesc
)

func file_proto_oidc_models_idp_proto_rawDescGZIP() []byte {
	file_proto_oidc_models_idp_proto_rawDescOnce.Do(func() {
		file_proto_oidc_models_idp_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_oidc_models_idp_proto_rawDescData)
	})
	return file_proto_oidc_models_idp_proto_rawDescData
}

var file_proto_oidc_models_idp_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_oidc_models_idp_proto_goTypes = []interface{}{
	(*OIDCProtocol)(nil),            // 0: proto.oidc.models.OIDCProtocol
	(*GithubOAuth2Protocol)(nil),    // 1: proto.oidc.models.GithubOAuth2Protocol
	(*OAuth2Protocol)(nil),          // 2: proto.oidc.models.OAuth2Protocol
	(*Protocol)(nil),                // 3: proto.oidc.models.Protocol
	(*ProtocolUpdate)(nil),          // 4: proto.oidc.models.ProtocolUpdate
	(*IDP)(nil),                     // 5: proto.oidc.models.IDP
	(*IDPs)(nil),                    // 6: proto.oidc.models.IDPs
	(*IDPUpdate)(nil),               // 7: proto.oidc.models.IDPUpdate
	nil,                             // 8: proto.oidc.models.IDP.MetadataEntry
	(*wrapperspb.BoolValue)(nil),    // 9: google.protobuf.BoolValue
	(*wrapperspb.StringValue)(nil),  // 10: google.protobuf.StringValue
	(*types.StringMapUpdate)(nil),   // 11: proto.types.StringMapUpdate
	(*types.StringArrayUpdate)(nil), // 12: proto.types.StringArrayUpdate
}
var file_proto_oidc_models_idp_proto_depIdxs = []int32{
	0,  // 0: proto.oidc.models.Protocol.oidc:type_name -> proto.oidc.models.OIDCProtocol
	2,  // 1: proto.oidc.models.Protocol.oauth2:type_name -> proto.oidc.models.OAuth2Protocol
	1,  // 2: proto.oidc.models.Protocol.github:type_name -> proto.oidc.models.GithubOAuth2Protocol
	3,  // 3: proto.oidc.models.ProtocolUpdate.value:type_name -> proto.oidc.models.Protocol
	3,  // 4: proto.oidc.models.IDP.protocol:type_name -> proto.oidc.models.Protocol
	8,  // 5: proto.oidc.models.IDP.metadata:type_name -> proto.oidc.models.IDP.MetadataEntry
	5,  // 6: proto.oidc.models.IDPs.idps:type_name -> proto.oidc.models.IDP
	9,  // 7: proto.oidc.models.IDPUpdate.enabled:type_name -> google.protobuf.BoolValue
	10, // 8: proto.oidc.models.IDPUpdate.slug:type_name -> google.protobuf.StringValue
	10, // 9: proto.oidc.models.IDPUpdate.name:type_name -> google.protobuf.StringValue
	10, // 10: proto.oidc.models.IDPUpdate.description:type_name -> google.protobuf.StringValue
	4,  // 11: proto.oidc.models.IDPUpdate.protocol:type_name -> proto.oidc.models.ProtocolUpdate
	11, // 12: proto.oidc.models.IDPUpdate.metadata:type_name -> proto.types.StringMapUpdate
	12, // 13: proto.oidc.models.IDPUpdate.claimed_domains:type_name -> proto.types.StringArrayUpdate
	9,  // 14: proto.oidc.models.IDPUpdate.hidden:type_name -> google.protobuf.BoolValue
	9,  // 15: proto.oidc.models.IDPUpdate.email_verification_required:type_name -> google.protobuf.BoolValue
	9,  // 16: proto.oidc.models.IDPUpdate.auto_create:type_name -> google.protobuf.BoolValue
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_proto_oidc_models_idp_proto_init() }
func file_proto_oidc_models_idp_proto_init() {
	if File_proto_oidc_models_idp_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_oidc_models_idp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OIDCProtocol); i {
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
		file_proto_oidc_models_idp_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GithubOAuth2Protocol); i {
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
		file_proto_oidc_models_idp_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OAuth2Protocol); i {
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
		file_proto_oidc_models_idp_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Protocol); i {
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
		file_proto_oidc_models_idp_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtocolUpdate); i {
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
		file_proto_oidc_models_idp_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IDP); i {
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
		file_proto_oidc_models_idp_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IDPs); i {
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
		file_proto_oidc_models_idp_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IDPUpdate); i {
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
	file_proto_oidc_models_idp_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Protocol_Oidc)(nil),
		(*Protocol_Oauth2)(nil),
		(*Protocol_Github)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_oidc_models_idp_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_oidc_models_idp_proto_goTypes,
		DependencyIndexes: file_proto_oidc_models_idp_proto_depIdxs,
		MessageInfos:      file_proto_oidc_models_idp_proto_msgTypes,
	}.Build()
	File_proto_oidc_models_idp_proto = out.File
	file_proto_oidc_models_idp_proto_rawDesc = nil
	file_proto_oidc_models_idp_proto_goTypes = nil
	file_proto_oidc_models_idp_proto_depIdxs = nil
}
