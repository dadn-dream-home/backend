// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.0
// source: protobuf/backend.proto

package protobuf

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

type FeedType int32

const (
	FeedType_TEMPERATURE FeedType = 0
	FeedType_HUMIDITY    FeedType = 1
	FeedType_LIGHT       FeedType = 2
)

// Enum value maps for FeedType.
var (
	FeedType_name = map[int32]string{
		0: "TEMPERATURE",
		1: "HUMIDITY",
		2: "LIGHT",
	}
	FeedType_value = map[string]int32{
		"TEMPERATURE": 0,
		"HUMIDITY":    1,
		"LIGHT":       2,
	}
)

func (x FeedType) Enum() *FeedType {
	p := new(FeedType)
	*p = x
	return p
}

func (x FeedType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FeedType) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_backend_proto_enumTypes[0].Descriptor()
}

func (FeedType) Type() protoreflect.EnumType {
	return &file_protobuf_backend_proto_enumTypes[0]
}

func (x FeedType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FeedType.Descriptor instead.
func (FeedType) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{0}
}

type StreamFeedValuesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *StreamFeedValuesRequest) Reset() {
	*x = StreamFeedValuesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamFeedValuesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamFeedValuesRequest) ProtoMessage() {}

func (x *StreamFeedValuesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamFeedValuesRequest.ProtoReflect.Descriptor instead.
func (*StreamFeedValuesRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{0}
}

func (x *StreamFeedValuesRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type StreamFeedValuesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float32 `protobuf:"fixed32,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StreamFeedValuesResponse) Reset() {
	*x = StreamFeedValuesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamFeedValuesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamFeedValuesResponse) ProtoMessage() {}

func (x *StreamFeedValuesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamFeedValuesResponse.ProtoReflect.Descriptor instead.
func (*StreamFeedValuesResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{1}
}

func (x *StreamFeedValuesResponse) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type StreamActuatorStatesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *StreamActuatorStatesRequest) Reset() {
	*x = StreamActuatorStatesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamActuatorStatesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamActuatorStatesRequest) ProtoMessage() {}

func (x *StreamActuatorStatesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamActuatorStatesRequest.ProtoReflect.Descriptor instead.
func (*StreamActuatorStatesRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{2}
}

func (x *StreamActuatorStatesRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type StreamActuatorStatesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State bool `protobuf:"varint,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *StreamActuatorStatesResponse) Reset() {
	*x = StreamActuatorStatesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamActuatorStatesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamActuatorStatesResponse) ProtoMessage() {}

func (x *StreamActuatorStatesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamActuatorStatesResponse.ProtoReflect.Descriptor instead.
func (*StreamActuatorStatesResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{3}
}

func (x *StreamActuatorStatesResponse) GetState() bool {
	if x != nil {
		return x.State
	}
	return false
}

type ListFeedsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListFeedsRequest) Reset() {
	*x = ListFeedsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFeedsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFeedsRequest) ProtoMessage() {}

func (x *ListFeedsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFeedsRequest.ProtoReflect.Descriptor instead.
func (*ListFeedsRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{4}
}

type ListFeedsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Feeds []*Feed `protobuf:"bytes,1,rep,name=feeds,proto3" json:"feeds,omitempty"`
}

func (x *ListFeedsResponse) Reset() {
	*x = ListFeedsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFeedsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFeedsResponse) ProtoMessage() {}

func (x *ListFeedsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFeedsResponse.ProtoReflect.Descriptor instead.
func (*ListFeedsResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{5}
}

func (x *ListFeedsResponse) GetFeeds() []*Feed {
	if x != nil {
		return x.Feeds
	}
	return nil
}

type CreateFeedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type FeedType `protobuf:"varint,2,opt,name=type,proto3,enum=protobuf.FeedType" json:"type,omitempty"`
}

func (x *CreateFeedRequest) Reset() {
	*x = CreateFeedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFeedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFeedRequest) ProtoMessage() {}

func (x *CreateFeedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFeedRequest.ProtoReflect.Descriptor instead.
func (*CreateFeedRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{6}
}

func (x *CreateFeedRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateFeedRequest) GetType() FeedType {
	if x != nil {
		return x.Type
	}
	return FeedType_TEMPERATURE
}

type CreateFeedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type FeedType `protobuf:"varint,2,opt,name=type,proto3,enum=protobuf.FeedType" json:"type,omitempty"`
}

func (x *CreateFeedResponse) Reset() {
	*x = CreateFeedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFeedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFeedResponse) ProtoMessage() {}

func (x *CreateFeedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFeedResponse.ProtoReflect.Descriptor instead.
func (*CreateFeedResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{7}
}

func (x *CreateFeedResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateFeedResponse) GetType() FeedType {
	if x != nil {
		return x.Type
	}
	return FeedType_TEMPERATURE
}

type DeleteFeedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteFeedRequest) Reset() {
	*x = DeleteFeedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFeedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFeedRequest) ProtoMessage() {}

func (x *DeleteFeedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFeedRequest.ProtoReflect.Descriptor instead.
func (*DeleteFeedRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteFeedRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteFeedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteFeedResponse) Reset() {
	*x = DeleteFeedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFeedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFeedResponse) ProtoMessage() {}

func (x *DeleteFeedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFeedResponse.ProtoReflect.Descriptor instead.
func (*DeleteFeedResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{9}
}

type Feed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Type        FeedType `protobuf:"varint,3,opt,name=type,proto3,enum=protobuf.FeedType" json:"type,omitempty"`
}

func (x *Feed) Reset() {
	*x = Feed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_backend_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Feed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Feed) ProtoMessage() {}

func (x *Feed) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_backend_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Feed.ProtoReflect.Descriptor instead.
func (*Feed) Descriptor() ([]byte, []int) {
	return file_protobuf_backend_proto_rawDescGZIP(), []int{10}
}

func (x *Feed) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Feed) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Feed) GetType() FeedType {
	if x != nil {
		return x.Type
	}
	return FeedType_TEMPERATURE
}

var File_protobuf_backend_proto protoreflect.FileDescriptor

var file_protobuf_backend_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x22, 0x29, 0x0a, 0x17, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46, 0x65, 0x65, 0x64,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x30, 0x0a,
	0x18, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46, 0x65, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0x2d, 0x0a, 0x1b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x41, 0x63, 0x74, 0x75, 0x61, 0x74, 0x6f,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x34,
	0x0a, 0x1c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x41, 0x63, 0x74, 0x75, 0x61, 0x74, 0x6f, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x65, 0x65, 0x64,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x39, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74,
	0x46, 0x65, 0x65, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a,
	0x05, 0x66, 0x65, 0x65, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x65, 0x65, 0x64, 0x52, 0x05, 0x66, 0x65,
	0x65, 0x64, 0x73, 0x22, 0x4b, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x65, 0x65, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x22, 0x4c, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x46, 0x65, 0x65, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x23,
	0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x14, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x65, 0x65,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x60, 0x0a, 0x04, 0x46, 0x65, 0x65,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x65, 0x65,
	0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x2a, 0x34, 0x0a, 0x08, 0x46,
	0x65, 0x65, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x45, 0x4d, 0x50, 0x45,
	0x52, 0x41, 0x54, 0x55, 0x52, 0x45, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x48, 0x55, 0x4d, 0x49,
	0x44, 0x49, 0x54, 0x59, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x49, 0x47, 0x48, 0x54, 0x10,
	0x02, 0x32, 0xb0, 0x03, 0x0a, 0x0e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x5d, 0x0a, 0x12, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65,
	0x6e, 0x73, 0x6f, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46, 0x65, 0x65, 0x64,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46,
	0x65, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x30, 0x01, 0x12, 0x67, 0x0a, 0x14, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x41, 0x63, 0x74,
	0x75, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x12, 0x25, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x41, 0x63, 0x74,
	0x75, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x41, 0x63, 0x74, 0x75, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x44, 0x0a, 0x09,
	0x4c, 0x69, 0x73, 0x74, 0x46, 0x65, 0x65, 0x64, 0x73, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x65, 0x65, 0x64, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x65, 0x65, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64,
	0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46,
	0x65, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x65, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x64, 0x6e, 0x2d, 0x64, 0x72, 0x65, 0x61, 0x6d, 0x2d, 0x68, 0x6f,
	0x6d, 0x65, 0x2f, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_backend_proto_rawDescOnce sync.Once
	file_protobuf_backend_proto_rawDescData = file_protobuf_backend_proto_rawDesc
)

func file_protobuf_backend_proto_rawDescGZIP() []byte {
	file_protobuf_backend_proto_rawDescOnce.Do(func() {
		file_protobuf_backend_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_backend_proto_rawDescData)
	})
	return file_protobuf_backend_proto_rawDescData
}

var file_protobuf_backend_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protobuf_backend_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_protobuf_backend_proto_goTypes = []interface{}{
	(FeedType)(0),                        // 0: protobuf.FeedType
	(*StreamFeedValuesRequest)(nil),      // 1: protobuf.StreamFeedValuesRequest
	(*StreamFeedValuesResponse)(nil),     // 2: protobuf.StreamFeedValuesResponse
	(*StreamActuatorStatesRequest)(nil),  // 3: protobuf.StreamActuatorStatesRequest
	(*StreamActuatorStatesResponse)(nil), // 4: protobuf.StreamActuatorStatesResponse
	(*ListFeedsRequest)(nil),             // 5: protobuf.ListFeedsRequest
	(*ListFeedsResponse)(nil),            // 6: protobuf.ListFeedsResponse
	(*CreateFeedRequest)(nil),            // 7: protobuf.CreateFeedRequest
	(*CreateFeedResponse)(nil),           // 8: protobuf.CreateFeedResponse
	(*DeleteFeedRequest)(nil),            // 9: protobuf.DeleteFeedRequest
	(*DeleteFeedResponse)(nil),           // 10: protobuf.DeleteFeedResponse
	(*Feed)(nil),                         // 11: protobuf.Feed
}
var file_protobuf_backend_proto_depIdxs = []int32{
	11, // 0: protobuf.ListFeedsResponse.feeds:type_name -> protobuf.Feed
	0,  // 1: protobuf.CreateFeedRequest.type:type_name -> protobuf.FeedType
	0,  // 2: protobuf.CreateFeedResponse.type:type_name -> protobuf.FeedType
	0,  // 3: protobuf.Feed.type:type_name -> protobuf.FeedType
	1,  // 4: protobuf.BackendService.StreamSensorValues:input_type -> protobuf.StreamFeedValuesRequest
	3,  // 5: protobuf.BackendService.StreamActuatorStates:input_type -> protobuf.StreamActuatorStatesRequest
	5,  // 6: protobuf.BackendService.ListFeeds:input_type -> protobuf.ListFeedsRequest
	7,  // 7: protobuf.BackendService.CreateFeed:input_type -> protobuf.CreateFeedRequest
	9,  // 8: protobuf.BackendService.DeleteFeed:input_type -> protobuf.DeleteFeedRequest
	2,  // 9: protobuf.BackendService.StreamSensorValues:output_type -> protobuf.StreamFeedValuesResponse
	4,  // 10: protobuf.BackendService.StreamActuatorStates:output_type -> protobuf.StreamActuatorStatesResponse
	6,  // 11: protobuf.BackendService.ListFeeds:output_type -> protobuf.ListFeedsResponse
	8,  // 12: protobuf.BackendService.CreateFeed:output_type -> protobuf.CreateFeedResponse
	10, // 13: protobuf.BackendService.DeleteFeed:output_type -> protobuf.DeleteFeedResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_protobuf_backend_proto_init() }
func file_protobuf_backend_proto_init() {
	if File_protobuf_backend_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_backend_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamFeedValuesRequest); i {
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
		file_protobuf_backend_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamFeedValuesResponse); i {
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
		file_protobuf_backend_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamActuatorStatesRequest); i {
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
		file_protobuf_backend_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamActuatorStatesResponse); i {
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
		file_protobuf_backend_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListFeedsRequest); i {
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
		file_protobuf_backend_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListFeedsResponse); i {
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
		file_protobuf_backend_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFeedRequest); i {
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
		file_protobuf_backend_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFeedResponse); i {
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
		file_protobuf_backend_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFeedRequest); i {
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
		file_protobuf_backend_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFeedResponse); i {
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
		file_protobuf_backend_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Feed); i {
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
			RawDescriptor: file_protobuf_backend_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_backend_proto_goTypes,
		DependencyIndexes: file_protobuf_backend_proto_depIdxs,
		EnumInfos:         file_protobuf_backend_proto_enumTypes,
		MessageInfos:      file_protobuf_backend_proto_msgTypes,
	}.Build()
	File_protobuf_backend_proto = out.File
	file_protobuf_backend_proto_rawDesc = nil
	file_protobuf_backend_proto_goTypes = nil
	file_protobuf_backend_proto_depIdxs = nil
}
