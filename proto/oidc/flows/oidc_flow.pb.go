// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: proto/oidc/flows/oidc_flow.proto

package flows

import (
	models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StoreAuthorizationRequestStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State              string                     `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	AuthorizationRequestState *models.AuthorizationRequestState `protobuf:"bytes,2,opt,name=authorization_final,json=authorizationFinal,proto3" json:"authorization_final,omitempty"`
}

func (x *StoreAuthorizationRequestStateRequest) Reset() {
	*x = StoreAuthorizationRequestStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreAuthorizationRequestStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreAuthorizationRequestStateRequest) ProtoMessage() {}

func (x *StoreAuthorizationRequestStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreAuthorizationRequestStateRequest.ProtoReflect.Descriptor instead.
func (*StoreAuthorizationRequestStateRequest) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{0}
}

func (x *StoreAuthorizationRequestStateRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *StoreAuthorizationRequestStateRequest) GetAuthorizationRequestState() *models.AuthorizationRequestState {
	if x != nil {
		return x.AuthorizationRequestState
	}
	return nil
}

type StoreAuthorizationRequestStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StoreAuthorizationRequestStateResponse) Reset() {
	*x = StoreAuthorizationRequestStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreAuthorizationRequestStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreAuthorizationRequestStateResponse) ProtoMessage() {}

func (x *StoreAuthorizationRequestStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreAuthorizationRequestStateResponse.ProtoReflect.Descriptor instead.
func (*StoreAuthorizationRequestStateResponse) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{1}
}

type GetAuthorizationRequestStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State string `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *GetAuthorizationRequestStateRequest) Reset() {
	*x = GetAuthorizationRequestStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthorizationRequestStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthorizationRequestStateRequest) ProtoMessage() {}

func (x *GetAuthorizationRequestStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthorizationRequestStateRequest.ProtoReflect.Descriptor instead.
func (*GetAuthorizationRequestStateRequest) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{2}
}

func (x *GetAuthorizationRequestStateRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type GetAuthorizationRequestStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthorizationRequestState *models.AuthorizationRequestState `protobuf:"bytes,1,opt,name=authorization_final,json=authorizationFinal,proto3" json:"authorization_final,omitempty"`
}

func (x *GetAuthorizationRequestStateResponse) Reset() {
	*x = GetAuthorizationRequestStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthorizationRequestStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthorizationRequestStateResponse) ProtoMessage() {}

func (x *GetAuthorizationRequestStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthorizationRequestStateResponse.ProtoReflect.Descriptor instead.
func (*GetAuthorizationRequestStateResponse) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{3}
}

func (x *GetAuthorizationRequestStateResponse) GetAuthorizationRequestState() *models.AuthorizationRequestState {
	if x != nil {
		return x.AuthorizationRequestState
	}
	return nil
}

type DeleteAuthorizationRequestStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State string `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *DeleteAuthorizationRequestStateRequest) Reset() {
	*x = DeleteAuthorizationRequestStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAuthorizationRequestStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAuthorizationRequestStateRequest) ProtoMessage() {}

func (x *DeleteAuthorizationRequestStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAuthorizationRequestStateRequest.ProtoReflect.Descriptor instead.
func (*DeleteAuthorizationRequestStateRequest) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteAuthorizationRequestStateRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type DeleteAuthorizationRequestStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteAuthorizationRequestStateResponse) Reset() {
	*x = DeleteAuthorizationRequestStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAuthorizationRequestStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAuthorizationRequestStateResponse) ProtoMessage() {}

func (x *DeleteAuthorizationRequestStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAuthorizationRequestStateResponse.ProtoReflect.Descriptor instead.
func (*DeleteAuthorizationRequestStateResponse) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{5}
}

type StoreExternalOauth2FinalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State               string                      `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	ExternalOauth2Final *models.ExternalOauth2Final `protobuf:"bytes,2,opt,name=external_oauth2_final,json=externalOauth2Final,proto3" json:"external_oauth2_final,omitempty"`
}

func (x *StoreExternalOauth2FinalRequest) Reset() {
	*x = StoreExternalOauth2FinalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreExternalOauth2FinalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreExternalOauth2FinalRequest) ProtoMessage() {}

func (x *StoreExternalOauth2FinalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreExternalOauth2FinalRequest.ProtoReflect.Descriptor instead.
func (*StoreExternalOauth2FinalRequest) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{6}
}

func (x *StoreExternalOauth2FinalRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *StoreExternalOauth2FinalRequest) GetExternalOauth2Final() *models.ExternalOauth2Final {
	if x != nil {
		return x.ExternalOauth2Final
	}
	return nil
}

type StoreExternalOauth2FinalResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StoreExternalOauth2FinalResponse) Reset() {
	*x = StoreExternalOauth2FinalResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreExternalOauth2FinalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreExternalOauth2FinalResponse) ProtoMessage() {}

func (x *StoreExternalOauth2FinalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreExternalOauth2FinalResponse.ProtoReflect.Descriptor instead.
func (*StoreExternalOauth2FinalResponse) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{7}
}

type GetExternalOauth2FinalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State string `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *GetExternalOauth2FinalRequest) Reset() {
	*x = GetExternalOauth2FinalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetExternalOauth2FinalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetExternalOauth2FinalRequest) ProtoMessage() {}

func (x *GetExternalOauth2FinalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetExternalOauth2FinalRequest.ProtoReflect.Descriptor instead.
func (*GetExternalOauth2FinalRequest) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{8}
}

func (x *GetExternalOauth2FinalRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type GetExternalOauth2FinalResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExternalOauth2Final *models.ExternalOauth2Final `protobuf:"bytes,1,opt,name=external_oauth2_final,json=externalOauth2Final,proto3" json:"external_oauth2_final,omitempty"`
}

func (x *GetExternalOauth2FinalResponse) Reset() {
	*x = GetExternalOauth2FinalResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetExternalOauth2FinalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetExternalOauth2FinalResponse) ProtoMessage() {}

func (x *GetExternalOauth2FinalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetExternalOauth2FinalResponse.ProtoReflect.Descriptor instead.
func (*GetExternalOauth2FinalResponse) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{9}
}

func (x *GetExternalOauth2FinalResponse) GetExternalOauth2Final() *models.ExternalOauth2Final {
	if x != nil {
		return x.ExternalOauth2Final
	}
	return nil
}

type DeleteExternalOauth2FinalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State string `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *DeleteExternalOauth2FinalRequest) Reset() {
	*x = DeleteExternalOauth2FinalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteExternalOauth2FinalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteExternalOauth2FinalRequest) ProtoMessage() {}

func (x *DeleteExternalOauth2FinalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteExternalOauth2FinalRequest.ProtoReflect.Descriptor instead.
func (*DeleteExternalOauth2FinalRequest) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteExternalOauth2FinalRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type DeleteExternalOauth2FinalResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteExternalOauth2FinalResponse) Reset() {
	*x = DeleteExternalOauth2FinalResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteExternalOauth2FinalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteExternalOauth2FinalResponse) ProtoMessage() {}

func (x *DeleteExternalOauth2FinalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oidc_flows_oidc_flow_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteExternalOauth2FinalResponse.ProtoReflect.Descriptor instead.
func (*DeleteExternalOauth2FinalResponse) Descriptor() ([]byte, []int) {
	return file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP(), []int{11}
}

var File_proto_oidc_flows_oidc_flow_proto protoreflect.FileDescriptor

var file_proto_oidc_flows_oidc_flow_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x69, 0x64, 0x63, 0x2f, 0x66, 0x6c, 0x6f,
	0x77, 0x73, 0x2f, 0x6f, 0x69, 0x64, 0x63, 0x5f, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x66,
	0x6c, 0x6f, 0x77, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x69,
	0x64, 0x63, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x01, 0x0a, 0x1e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x56, 0x0a,
	0x13, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66,
	0x69, 0x6e, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x52, 0x12, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x22, 0x21, 0x0a, 0x1f, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x34, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x77,
	0x0a, 0x1d, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x56, 0x0a, 0x13, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69,
	0x6e, 0x61, 0x6c, 0x52, 0x12, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x22, 0x37, 0x0a, 0x1f, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69,
	0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x22, 0x22, 0x0a, 0x20, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x93, 0x01, 0x0a, 0x1f, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x5a,
	0x0a, 0x15, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x6f, 0x61, 0x75, 0x74, 0x68,
	0x32, 0x5f, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x13, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f,
	0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x22, 0x22, 0x0a, 0x20, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68,
	0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x35,
	0x0a, 0x1d, 0x47, 0x65, 0x74, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75,
	0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x7c, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x45, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x15, 0x65, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x5f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x5f, 0x66, 0x69, 0x6e, 0x61, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f,
	0x69, 0x64, 0x63, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x13,
	0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69,
	0x6e, 0x61, 0x6c, 0x22, 0x38, 0x0a, 0x20, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x23, 0x0a,
	0x21, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f,
	0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x32, 0x94, 0x03, 0x0a, 0x0d, 0x4f, 0x49, 0x44, 0x43, 0x46, 0x6c, 0x6f, 0x77, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x12, 0x80, 0x01, 0x0a, 0x17, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x12, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x7a, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x12, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x83, 0x01, 0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x12, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63,
	0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0xa7, 0x03, 0x0a, 0x17, 0x45, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x6c, 0x6f, 0x77,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x18, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e,
	0x61, 0x6c, 0x12, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69,
	0x64, 0x63, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x7d, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x2f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69,
	0x64, 0x63, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f,
	0x69, 0x64, 0x63, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x78, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x86, 0x01, 0x0a, 0x19, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75,
	0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32,
	0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6f, 0x69, 0x64, 0x63, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4f, 0x61,
	0x75, 0x74, 0x68, 0x32, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x8d, 0x01, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x66,
	0x66, 0x79, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x72, 0x61, 0x67, 0x65, 0x6f, 0x69, 0x64, 0x63,
	0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x50, 0x01, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x2d, 0x62, 0x75, 0x6e, 0x6e,
	0x79, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x63, 0x6f, 0x72, 0x65, 0x2d, 0x72, 0x61, 0x67,
	0x65, 0x2d, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6f, 0x69, 0x64, 0x63, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x3b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0xaa, 0x02, 0x1f, 0x46, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x42, 0x75, 0x6e, 0x6e, 0x79, 0x2e,
	0x52, 0x61, 0x67, 0x65, 0x4f, 0x69, 0x64, 0x63, 0x2e, 0x4f, 0x69, 0x64, 0x63, 0x2e, 0x46, 0x6c,
	0x6f, 0x77, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_oidc_flows_oidc_flow_proto_rawDescOnce sync.Once
	file_proto_oidc_flows_oidc_flow_proto_rawDescData = file_proto_oidc_flows_oidc_flow_proto_rawDesc
)

func file_proto_oidc_flows_oidc_flow_proto_rawDescGZIP() []byte {
	file_proto_oidc_flows_oidc_flow_proto_rawDescOnce.Do(func() {
		file_proto_oidc_flows_oidc_flow_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_oidc_flows_oidc_flow_proto_rawDescData)
	})
	return file_proto_oidc_flows_oidc_flow_proto_rawDescData
}

var file_proto_oidc_flows_oidc_flow_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_proto_oidc_flows_oidc_flow_proto_goTypes = []interface{}{
	(*StoreAuthorizationRequestStateRequest)(nil),    // 0: proto.oidc.flows.StoreAuthorizationRequestStateRequest
	(*StoreAuthorizationRequestStateResponse)(nil),   // 1: proto.oidc.flows.StoreAuthorizationRequestStateResponse
	(*GetAuthorizationRequestStateRequest)(nil),      // 2: proto.oidc.flows.GetAuthorizationRequestStateRequest
	(*GetAuthorizationRequestStateResponse)(nil),     // 3: proto.oidc.flows.GetAuthorizationRequestStateResponse
	(*DeleteAuthorizationRequestStateRequest)(nil),   // 4: proto.oidc.flows.DeleteAuthorizationRequestStateRequest
	(*DeleteAuthorizationRequestStateResponse)(nil),  // 5: proto.oidc.flows.DeleteAuthorizationRequestStateResponse
	(*StoreExternalOauth2FinalRequest)(nil),   // 6: proto.oidc.flows.StoreExternalOauth2FinalRequest
	(*StoreExternalOauth2FinalResponse)(nil),  // 7: proto.oidc.flows.StoreExternalOauth2FinalResponse
	(*GetExternalOauth2FinalRequest)(nil),     // 8: proto.oidc.flows.GetExternalOauth2FinalRequest
	(*GetExternalOauth2FinalResponse)(nil),    // 9: proto.oidc.flows.GetExternalOauth2FinalResponse
	(*DeleteExternalOauth2FinalRequest)(nil),  // 10: proto.oidc.flows.DeleteExternalOauth2FinalRequest
	(*DeleteExternalOauth2FinalResponse)(nil), // 11: proto.oidc.flows.DeleteExternalOauth2FinalResponse
	(*models.AuthorizationRequestState)(nil),         // 12: proto.oidc.models.AuthorizationRequestState
	(*models.ExternalOauth2Final)(nil),        // 13: proto.oidc.models.ExternalOauth2Final
}
var file_proto_oidc_flows_oidc_flow_proto_depIdxs = []int32{
	12, // 0: proto.oidc.flows.StoreAuthorizationRequestStateRequest.authorization_final:type_name -> proto.oidc.models.AuthorizationRequestState
	12, // 1: proto.oidc.flows.GetAuthorizationRequestStateResponse.authorization_final:type_name -> proto.oidc.models.AuthorizationRequestState
	13, // 2: proto.oidc.flows.StoreExternalOauth2FinalRequest.external_oauth2_final:type_name -> proto.oidc.models.ExternalOauth2Final
	13, // 3: proto.oidc.flows.GetExternalOauth2FinalResponse.external_oauth2_final:type_name -> proto.oidc.models.ExternalOauth2Final
	0,  // 4: proto.oidc.flows.OIDCFlowStore.StoreAuthorizationRequestState:input_type -> proto.oidc.flows.StoreAuthorizationRequestStateRequest
	2,  // 5: proto.oidc.flows.OIDCFlowStore.GetAuthorizationRequestState:input_type -> proto.oidc.flows.GetAuthorizationRequestStateRequest
	4,  // 6: proto.oidc.flows.OIDCFlowStore.DeleteAuthorizationRequestState:input_type -> proto.oidc.flows.DeleteAuthorizationRequestStateRequest
	6,  // 7: proto.oidc.flows.ExternalOauth2FlowStore.StoreExternalOauth2Final:input_type -> proto.oidc.flows.StoreExternalOauth2FinalRequest
	8,  // 8: proto.oidc.flows.ExternalOauth2FlowStore.GetExternalOauth2Final:input_type -> proto.oidc.flows.GetExternalOauth2FinalRequest
	10, // 9: proto.oidc.flows.ExternalOauth2FlowStore.DeleteExternalOauth2Final:input_type -> proto.oidc.flows.DeleteExternalOauth2FinalRequest
	1,  // 10: proto.oidc.flows.OIDCFlowStore.StoreAuthorizationRequestState:output_type -> proto.oidc.flows.StoreAuthorizationRequestStateResponse
	3,  // 11: proto.oidc.flows.OIDCFlowStore.GetAuthorizationRequestState:output_type -> proto.oidc.flows.GetAuthorizationRequestStateResponse
	5,  // 12: proto.oidc.flows.OIDCFlowStore.DeleteAuthorizationRequestState:output_type -> proto.oidc.flows.DeleteAuthorizationRequestStateResponse
	7,  // 13: proto.oidc.flows.ExternalOauth2FlowStore.StoreExternalOauth2Final:output_type -> proto.oidc.flows.StoreExternalOauth2FinalResponse
	9,  // 14: proto.oidc.flows.ExternalOauth2FlowStore.GetExternalOauth2Final:output_type -> proto.oidc.flows.GetExternalOauth2FinalResponse
	11, // 15: proto.oidc.flows.ExternalOauth2FlowStore.DeleteExternalOauth2Final:output_type -> proto.oidc.flows.DeleteExternalOauth2FinalResponse
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_oidc_flows_oidc_flow_proto_init() }
func file_proto_oidc_flows_oidc_flow_proto_init() {
	if File_proto_oidc_flows_oidc_flow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreAuthorizationRequestStateRequest); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreAuthorizationRequestStateResponse); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthorizationRequestStateRequest); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthorizationRequestStateResponse); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAuthorizationRequestStateRequest); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAuthorizationRequestStateResponse); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreExternalOauth2FinalRequest); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreExternalOauth2FinalResponse); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetExternalOauth2FinalRequest); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetExternalOauth2FinalResponse); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteExternalOauth2FinalRequest); i {
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
		file_proto_oidc_flows_oidc_flow_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteExternalOauth2FinalResponse); i {
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
			RawDescriptor: file_proto_oidc_flows_oidc_flow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_proto_oidc_flows_oidc_flow_proto_goTypes,
		DependencyIndexes: file_proto_oidc_flows_oidc_flow_proto_depIdxs,
		MessageInfos:      file_proto_oidc_flows_oidc_flow_proto_msgTypes,
	}.Build()
	File_proto_oidc_flows_oidc_flow_proto = out.File
	file_proto_oidc_flows_oidc_flow_proto_rawDesc = nil
	file_proto_oidc_flows_oidc_flow_proto_goTypes = nil
	file_proto_oidc_flows_oidc_flow_proto_depIdxs = nil
}