// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.1
// source: proto/types/pagination.proto

package types

import (
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

type Order int32

const (
	// Conventional default for enums. Do not use this.
	Order_ORDER_UNSPECIFIED Order = 0
	// Ascending order
	Order_ASC Order = 1
	// Descending order
	Order_DESC Order = 2
)

// Enum value maps for Order.
var (
	Order_name = map[int32]string{
		0: "ORDER_UNSPECIFIED",
		1: "ASC",
		2: "DESC",
	}
	Order_value = map[string]int32{
		"ORDER_UNSPECIFIED": 0,
		"ASC":               1,
		"DESC":              2,
	}
)

func (x Order) Enum() *Order {
	p := new(Order)
	*p = x
	return p
}

func (x Order) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Order) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_types_pagination_proto_enumTypes[0].Descriptor()
}

func (Order) Type() protoreflect.EnumType {
	return &file_proto_types_pagination_proto_enumTypes[0]
}

func (x Order) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Order.Descriptor instead.
func (Order) EnumDescriptor() ([]byte, []int) {
	return file_proto_types_pagination_proto_rawDescGZIP(), []int{0}
}

type Pagination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Maximum number of entries to retrieve
	Limit    uint32 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Iterator string `protobuf:"bytes,2,opt,name=iterator,proto3" json:"iterator,omitempty"`
	Order    Order  `protobuf:"varint,3,opt,name=order,proto3,enum=proto.types.Order" json:"order,omitempty"`
}

func (x *Pagination) Reset() {
	*x = Pagination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_pagination_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pagination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pagination) ProtoMessage() {}

func (x *Pagination) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_pagination_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pagination.ProtoReflect.Descriptor instead.
func (*Pagination) Descriptor() ([]byte, []int) {
	return file_proto_types_pagination_proto_rawDescGZIP(), []int{0}
}

func (x *Pagination) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *Pagination) GetIterator() string {
	if x != nil {
		return x.Iterator
	}
	return ""
}

func (x *Pagination) GetOrder() Order {
	if x != nil {
		return x.Order
	}
	return Order_ORDER_UNSPECIFIED
}

// PaginationResponse ...
type PaginationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Iterator     string `protobuf:"bytes,1,opt,name=iterator,proto3" json:"iterator,omitempty"`
	PrevIterator string `protobuf:"bytes,2,opt,name=prev_iterator,json=prevIterator,proto3" json:"prev_iterator,omitempty"`
	// Is a total count of records available?
	TotalAvailable bool `protobuf:"varint,3,opt,name=total_available,json=totalAvailable,proto3" json:"total_available,omitempty"`
	// Total number of records available (if totalAvailable = true)
	Total uint64 `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`
	// There is no more data
	NoMoreData bool `protobuf:"varint,5,opt,name=no_more_data,json=noMoreData,proto3" json:"no_more_data,omitempty"`
}

func (x *PaginationResponse) Reset() {
	*x = PaginationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_pagination_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaginationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationResponse) ProtoMessage() {}

func (x *PaginationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_pagination_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaginationResponse.ProtoReflect.Descriptor instead.
func (*PaginationResponse) Descriptor() ([]byte, []int) {
	return file_proto_types_pagination_proto_rawDescGZIP(), []int{1}
}

func (x *PaginationResponse) GetIterator() string {
	if x != nil {
		return x.Iterator
	}
	return ""
}

func (x *PaginationResponse) GetPrevIterator() string {
	if x != nil {
		return x.PrevIterator
	}
	return ""
}

func (x *PaginationResponse) GetTotalAvailable() bool {
	if x != nil {
		return x.TotalAvailable
	}
	return false
}

func (x *PaginationResponse) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *PaginationResponse) GetNoMoreData() bool {
	if x != nil {
		return x.NoMoreData
	}
	return false
}

var File_proto_types_pagination_proto protoreflect.FileDescriptor

var file_proto_types_pagination_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x68, 0x0a,
	0x0a, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x28, 0x0a,
	0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x22, 0xb6, 0x01, 0x0a, 0x12, 0x50, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72,
	0x65, 0x76, 0x5f, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x49, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12,
	0x27, 0x0a, 0x0f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x41,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x20,
	0x0a, 0x0c, 0x6e, 0x6f, 0x5f, 0x6d, 0x6f, 0x72, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x6e, 0x6f, 0x4d, 0x6f, 0x72, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x2a, 0x31, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x15, 0x0a, 0x11, 0x4f, 0x52, 0x44,
	0x45, 0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x07, 0x0a, 0x03, 0x41, 0x53, 0x43, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x45, 0x53,
	0x43, 0x10, 0x02, 0x42, 0x83, 0x01, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x66,
	0x66, 0x79, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x72, 0x61, 0x67, 0x65, 0x6f, 0x69, 0x64, 0x63,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x2d, 0x62, 0x75, 0x6e, 0x6e,
	0x79, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x63, 0x6f, 0x72, 0x65, 0x2d, 0x72, 0x61, 0x67,
	0x65, 0x2d, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x73, 0xaa, 0x02, 0x1a, 0x46,
	0x6c, 0x75, 0x66, 0x66, 0x79, 0x42, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x52, 0x61, 0x67, 0x65, 0x4f,
	0x69, 0x64, 0x63, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_types_pagination_proto_rawDescOnce sync.Once
	file_proto_types_pagination_proto_rawDescData = file_proto_types_pagination_proto_rawDesc
)

func file_proto_types_pagination_proto_rawDescGZIP() []byte {
	file_proto_types_pagination_proto_rawDescOnce.Do(func() {
		file_proto_types_pagination_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_types_pagination_proto_rawDescData)
	})
	return file_proto_types_pagination_proto_rawDescData
}

var file_proto_types_pagination_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_types_pagination_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_types_pagination_proto_goTypes = []any{
	(Order)(0),                 // 0: proto.types.Order
	(*Pagination)(nil),         // 1: proto.types.Pagination
	(*PaginationResponse)(nil), // 2: proto.types.PaginationResponse
}
var file_proto_types_pagination_proto_depIdxs = []int32{
	0, // 0: proto.types.Pagination.order:type_name -> proto.types.Order
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_types_pagination_proto_init() }
func file_proto_types_pagination_proto_init() {
	if File_proto_types_pagination_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_types_pagination_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Pagination); i {
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
		file_proto_types_pagination_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*PaginationResponse); i {
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
			RawDescriptor: file_proto_types_pagination_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_types_pagination_proto_goTypes,
		DependencyIndexes: file_proto_types_pagination_proto_depIdxs,
		EnumInfos:         file_proto_types_pagination_proto_enumTypes,
		MessageInfos:      file_proto_types_pagination_proto_msgTypes,
	}.Build()
	File_proto_types_pagination_proto = out.File
	file_proto_types_pagination_proto_rawDesc = nil
	file_proto_types_pagination_proto_goTypes = nil
	file_proto_types_pagination_proto_depIdxs = nil
}
