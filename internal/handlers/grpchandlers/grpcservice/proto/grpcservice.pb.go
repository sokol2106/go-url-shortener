// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.0--rc2
// source: grpcservice.proto

package proto

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

type CreateShortURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *CreateShortURLRequest) Reset() {
	*x = CreateShortURLRequest{}
	mi := &file_grpcservice_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateShortURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShortURLRequest) ProtoMessage() {}

func (x *CreateShortURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShortURLRequest.ProtoReflect.Descriptor instead.
func (*CreateShortURLRequest) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{0}
}

func (x *CreateShortURLRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CreateShortURLRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type CreateShortURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *CreateShortURLResponse) Reset() {
	*x = CreateShortURLResponse{}
	mi := &file_grpcservice_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateShortURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShortURLResponse) ProtoMessage() {}

func (x *CreateShortURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShortURLResponse.ProtoReflect.Descriptor instead.
func (*CreateShortURLResponse) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{1}
}

func (x *CreateShortURLResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type CreateShortURLRequestJSON struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *CreateShortURLRequestJSON) Reset() {
	*x = CreateShortURLRequestJSON{}
	mi := &file_grpcservice_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateShortURLRequestJSON) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShortURLRequestJSON) ProtoMessage() {}

func (x *CreateShortURLRequestJSON) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShortURLRequestJSON.ProtoReflect.Descriptor instead.
func (*CreateShortURLRequestJSON) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{2}
}

func (x *CreateShortURLRequestJSON) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CreateShortURLRequestJSON) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type CreateShortURLResponseJSON struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *CreateShortURLResponseJSON) Reset() {
	*x = CreateShortURLResponseJSON{}
	mi := &file_grpcservice_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateShortURLResponseJSON) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShortURLResponseJSON) ProtoMessage() {}

func (x *CreateShortURLResponseJSON) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShortURLResponseJSON.ProtoReflect.Descriptor instead.
func (*CreateShortURLResponseJSON) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{3}
}

func (x *CreateShortURLResponseJSON) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type GetOriginalURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *GetOriginalURLRequest) Reset() {
	*x = GetOriginalURLRequest{}
	mi := &file_grpcservice_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOriginalURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOriginalURLRequest) ProtoMessage() {}

func (x *GetOriginalURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOriginalURLRequest.ProtoReflect.Descriptor instead.
func (*GetOriginalURLRequest) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{4}
}

func (x *GetOriginalURLRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type GetOriginalURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalUrl string `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
}

func (x *GetOriginalURLResponse) Reset() {
	*x = GetOriginalURLResponse{}
	mi := &file_grpcservice_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOriginalURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOriginalURLResponse) ProtoMessage() {}

func (x *GetOriginalURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOriginalURLResponse.ProtoReflect.Descriptor instead.
func (*GetOriginalURLResponse) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{5}
}

func (x *GetOriginalURLResponse) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

type GetUserShortenedURLsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *GetUserShortenedURLsRequest) Reset() {
	*x = GetUserShortenedURLsRequest{}
	mi := &file_grpcservice_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserShortenedURLsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserShortenedURLsRequest) ProtoMessage() {}

func (x *GetUserShortenedURLsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserShortenedURLsRequest.ProtoReflect.Descriptor instead.
func (*GetUserShortenedURLsRequest) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{6}
}

func (x *GetUserShortenedURLsRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type GetUserShortenedURLsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []string `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *GetUserShortenedURLsResponse) Reset() {
	*x = GetUserShortenedURLsResponse{}
	mi := &file_grpcservice_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserShortenedURLsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserShortenedURLsResponse) ProtoMessage() {}

func (x *GetUserShortenedURLsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserShortenedURLsResponse.ProtoReflect.Descriptor instead.
func (*GetUserShortenedURLsResponse) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{7}
}

func (x *GetUserShortenedURLsResponse) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

type DeleteUserShortenedURLsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Urls  []string `protobuf:"bytes,2,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *DeleteUserShortenedURLsRequest) Reset() {
	*x = DeleteUserShortenedURLsRequest{}
	mi := &file_grpcservice_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserShortenedURLsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserShortenedURLsRequest) ProtoMessage() {}

func (x *DeleteUserShortenedURLsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserShortenedURLsRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserShortenedURLsRequest) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteUserShortenedURLsRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *DeleteUserShortenedURLsRequest) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

type DeleteUserShortenedURLsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteUserShortenedURLsResponse) Reset() {
	*x = DeleteUserShortenedURLsResponse{}
	mi := &file_grpcservice_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserShortenedURLsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserShortenedURLsResponse) ProtoMessage() {}

func (x *DeleteUserShortenedURLsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpcservice_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserShortenedURLsResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserShortenedURLsResponse) Descriptor() ([]byte, []int) {
	return file_grpcservice_proto_rawDescGZIP(), []int{9}
}

var File_grpcservice_proto protoreflect.FileDescriptor

var file_grpcservice_proto_rawDesc = []byte{
	0x0a, 0x11, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x22, 0x3f, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55,
	0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x30, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x43, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4a, 0x53, 0x4f, 0x4e,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x34, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x4a, 0x53, 0x4f, 0x4e, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x2b,
	0x0a, 0x15, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x3b, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61,
	0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c, 0x22, 0x33, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x32, 0x0a,
	0x1c, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65,
	0x64, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x75, 0x72, 0x6c,
	0x73, 0x22, 0x4a, 0x0a, 0x1e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72, 0x6c,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x21, 0x0a,
	0x1f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0x8e, 0x04, 0x0a, 0x0c, 0x55, 0x52, 0x4c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65,
	0x72, 0x12, 0x59, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x55, 0x52, 0x4c, 0x12, 0x22, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x12,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x4a, 0x53,
	0x4f, 0x4e, 0x12, 0x26, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4a, 0x53, 0x4f, 0x4e, 0x1a, 0x27, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4a,
	0x53, 0x4f, 0x4e, 0x12, 0x59, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x61, 0x6c, 0x55, 0x52, 0x4c, 0x12, 0x22, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55,
	0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6b,
	0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x64, 0x55, 0x52, 0x4c, 0x73, 0x12, 0x28, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x29, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x64, 0x55,
	0x52, 0x4c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x74, 0x0a, 0x17, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x64, 0x55, 0x52, 0x4c, 0x73, 0x12, 0x2b, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x13, 0x5a, 0x11, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpcservice_proto_rawDescOnce sync.Once
	file_grpcservice_proto_rawDescData = file_grpcservice_proto_rawDesc
)

func file_grpcservice_proto_rawDescGZIP() []byte {
	file_grpcservice_proto_rawDescOnce.Do(func() {
		file_grpcservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpcservice_proto_rawDescData)
	})
	return file_grpcservice_proto_rawDescData
}

var file_grpcservice_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_grpcservice_proto_goTypes = []any{
	(*CreateShortURLRequest)(nil),           // 0: grpcservice.CreateShortURLRequest
	(*CreateShortURLResponse)(nil),          // 1: grpcservice.CreateShortURLResponse
	(*CreateShortURLRequestJSON)(nil),       // 2: grpcservice.CreateShortURLRequestJSON
	(*CreateShortURLResponseJSON)(nil),      // 3: grpcservice.CreateShortURLResponseJSON
	(*GetOriginalURLRequest)(nil),           // 4: grpcservice.GetOriginalURLRequest
	(*GetOriginalURLResponse)(nil),          // 5: grpcservice.GetOriginalURLResponse
	(*GetUserShortenedURLsRequest)(nil),     // 6: grpcservice.GetUserShortenedURLsRequest
	(*GetUserShortenedURLsResponse)(nil),    // 7: grpcservice.GetUserShortenedURLsResponse
	(*DeleteUserShortenedURLsRequest)(nil),  // 8: grpcservice.DeleteUserShortenedURLsRequest
	(*DeleteUserShortenedURLsResponse)(nil), // 9: grpcservice.DeleteUserShortenedURLsResponse
}
var file_grpcservice_proto_depIdxs = []int32{
	0, // 0: grpcservice.URLShortener.CreateShortURL:input_type -> grpcservice.CreateShortURLRequest
	2, // 1: grpcservice.URLShortener.CreateShortURLJSON:input_type -> grpcservice.CreateShortURLRequestJSON
	4, // 2: grpcservice.URLShortener.GetOriginalURL:input_type -> grpcservice.GetOriginalURLRequest
	6, // 3: grpcservice.URLShortener.GetUserShortenedURLs:input_type -> grpcservice.GetUserShortenedURLsRequest
	8, // 4: grpcservice.URLShortener.DeleteUserShortenedURLs:input_type -> grpcservice.DeleteUserShortenedURLsRequest
	1, // 5: grpcservice.URLShortener.CreateShortURL:output_type -> grpcservice.CreateShortURLResponse
	3, // 6: grpcservice.URLShortener.CreateShortURLJSON:output_type -> grpcservice.CreateShortURLResponseJSON
	5, // 7: grpcservice.URLShortener.GetOriginalURL:output_type -> grpcservice.GetOriginalURLResponse
	7, // 8: grpcservice.URLShortener.GetUserShortenedURLs:output_type -> grpcservice.GetUserShortenedURLsResponse
	9, // 9: grpcservice.URLShortener.DeleteUserShortenedURLs:output_type -> grpcservice.DeleteUserShortenedURLsResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpcservice_proto_init() }
func file_grpcservice_proto_init() {
	if File_grpcservice_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpcservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpcservice_proto_goTypes,
		DependencyIndexes: file_grpcservice_proto_depIdxs,
		MessageInfos:      file_grpcservice_proto_msgTypes,
	}.Build()
	File_grpcservice_proto = out.File
	file_grpcservice_proto_rawDesc = nil
	file_grpcservice_proto_goTypes = nil
	file_grpcservice_proto_depIdxs = nil
}