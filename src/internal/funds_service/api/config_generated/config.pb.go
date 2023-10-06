// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: src/proto/funds_service/config.proto

package config_generated

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Chain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Chain) Reset() {
	*x = Chain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_funds_service_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chain) ProtoMessage() {}

func (x *Chain) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_funds_service_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chain.ProtoReflect.Descriptor instead.
func (*Chain) Descriptor() ([]byte, []int) {
	return file_src_proto_funds_service_config_proto_rawDescGZIP(), []int{0}
}

type SetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chains                    map[string]*anypb.Any `protobuf:"bytes,1,rep,name=chains,proto3" json:"chains,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ExpireTime                *int64                `protobuf:"varint,2,opt,name=expireTime,proto3,oneof" json:"expireTime,omitempty"`
	Mnemonic                  *string               `protobuf:"bytes,3,opt,name=mnemonic,proto3,oneof" json:"mnemonic,omitempty"`
	WalletCollectionThreshold *int64                `protobuf:"varint,4,opt,name=walletCollectionThreshold,proto3,oneof" json:"walletCollectionThreshold,omitempty"`
	MinGasThreshold           *int64                `protobuf:"varint,5,opt,name=minGasThreshold,proto3,oneof" json:"minGasThreshold,omitempty"`
	TransferGasAmount         *int64                `protobuf:"varint,6,opt,name=transferGasAmount,proto3,oneof" json:"transferGasAmount,omitempty"`
}

func (x *SetRequest) Reset() {
	*x = SetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_funds_service_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetRequest) ProtoMessage() {}

func (x *SetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_funds_service_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetRequest.ProtoReflect.Descriptor instead.
func (*SetRequest) Descriptor() ([]byte, []int) {
	return file_src_proto_funds_service_config_proto_rawDescGZIP(), []int{1}
}

func (x *SetRequest) GetChains() map[string]*anypb.Any {
	if x != nil {
		return x.Chains
	}
	return nil
}

func (x *SetRequest) GetExpireTime() int64 {
	if x != nil && x.ExpireTime != nil {
		return *x.ExpireTime
	}
	return 0
}

func (x *SetRequest) GetMnemonic() string {
	if x != nil && x.Mnemonic != nil {
		return *x.Mnemonic
	}
	return ""
}

func (x *SetRequest) GetWalletCollectionThreshold() int64 {
	if x != nil && x.WalletCollectionThreshold != nil {
		return *x.WalletCollectionThreshold
	}
	return 0
}

func (x *SetRequest) GetMinGasThreshold() int64 {
	if x != nil && x.MinGasThreshold != nil {
		return *x.MinGasThreshold
	}
	return 0
}

func (x *SetRequest) GetTransferGasAmount() int64 {
	if x != nil && x.TransferGasAmount != nil {
		return *x.TransferGasAmount
	}
	return 0
}

type LoadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chains                    map[string]*anypb.Any `protobuf:"bytes,1,rep,name=chains,proto3" json:"chains,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ExpireTime                int64                 `protobuf:"varint,2,opt,name=expireTime,proto3" json:"expireTime,omitempty"`
	Mnemonic                  string                `protobuf:"bytes,3,opt,name=mnemonic,proto3" json:"mnemonic,omitempty"`
	WalletCollectionThreshold float64               `protobuf:"fixed64,4,opt,name=walletCollectionThreshold,proto3" json:"walletCollectionThreshold,omitempty"`
	MinGasThreshold           float64               `protobuf:"fixed64,5,opt,name=minGasThreshold,proto3" json:"minGasThreshold,omitempty"`
	TransferGasAmount         float64               `protobuf:"fixed64,6,opt,name=transferGasAmount,proto3" json:"transferGasAmount,omitempty"`
}

func (x *LoadResponse) Reset() {
	*x = LoadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_funds_service_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadResponse) ProtoMessage() {}

func (x *LoadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_funds_service_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadResponse.ProtoReflect.Descriptor instead.
func (*LoadResponse) Descriptor() ([]byte, []int) {
	return file_src_proto_funds_service_config_proto_rawDescGZIP(), []int{2}
}

func (x *LoadResponse) GetChains() map[string]*anypb.Any {
	if x != nil {
		return x.Chains
	}
	return nil
}

func (x *LoadResponse) GetExpireTime() int64 {
	if x != nil {
		return x.ExpireTime
	}
	return 0
}

func (x *LoadResponse) GetMnemonic() string {
	if x != nil {
		return x.Mnemonic
	}
	return ""
}

func (x *LoadResponse) GetWalletCollectionThreshold() float64 {
	if x != nil {
		return x.WalletCollectionThreshold
	}
	return 0
}

func (x *LoadResponse) GetMinGasThreshold() float64 {
	if x != nil {
		return x.MinGasThreshold
	}
	return 0
}

func (x *LoadResponse) GetTransferGasAmount() float64 {
	if x != nil {
		return x.TransferGasAmount
	}
	return 0
}

type Chain_Tron struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RpcNodes  []string `protobuf:"bytes,1,rep,name=rpcNodes,proto3" json:"rpcNodes,omitempty"`
	HttpNodes []string `protobuf:"bytes,2,rep,name=httpNodes,proto3" json:"httpNodes,omitempty"`
	ApiKeys   []string `protobuf:"bytes,3,rep,name=apiKeys,proto3" json:"apiKeys,omitempty"`
}

func (x *Chain_Tron) Reset() {
	*x = Chain_Tron{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_funds_service_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chain_Tron) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chain_Tron) ProtoMessage() {}

func (x *Chain_Tron) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_funds_service_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chain_Tron.ProtoReflect.Descriptor instead.
func (*Chain_Tron) Descriptor() ([]byte, []int) {
	return file_src_proto_funds_service_config_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Chain_Tron) GetRpcNodes() []string {
	if x != nil {
		return x.RpcNodes
	}
	return nil
}

func (x *Chain_Tron) GetHttpNodes() []string {
	if x != nil {
		return x.HttpNodes
	}
	return nil
}

func (x *Chain_Tron) GetApiKeys() []string {
	if x != nil {
		return x.ApiKeys
	}
	return nil
}

var File_src_proto_funds_service_config_proto protoreflect.FileDescriptor

var file_src_proto_funds_service_config_proto_rawDesc = []byte{
	0x0a, 0x24, 0x73, 0x72, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x64,
	0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63,
	0x0a, 0x05, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x1a, 0x5a, 0x0a, 0x04, 0x54, 0x72, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x72, 0x70, 0x63, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x08, 0x72, 0x70, 0x63, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x68,
	0x74, 0x74, 0x70, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09,
	0x68, 0x74, 0x74, 0x70, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x70, 0x69,
	0x4b, 0x65, 0x79, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x61, 0x70, 0x69, 0x4b,
	0x65, 0x79, 0x73, 0x22, 0xdd, 0x03, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2f, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x73, 0x12, 0x23, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x6d, 0x6e, 0x65, 0x6d,
	0x6f, 0x6e, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x6d, 0x6e,
	0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x88, 0x01, 0x01, 0x12, 0x41, 0x0a, 0x19, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x68, 0x72,
	0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x19,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x0f,
	0x6d, 0x69, 0x6e, 0x47, 0x61, 0x73, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x03, 0x52, 0x0f, 0x6d, 0x69, 0x6e, 0x47, 0x61, 0x73, 0x54,
	0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x88, 0x01, 0x01, 0x12, 0x31, 0x0a, 0x11, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x47, 0x61, 0x73, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x48, 0x04, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x47, 0x61, 0x73, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x1a, 0x4f,
	0x0a, 0x0b, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42,
	0x0d, 0x0a, 0x0b, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x0b,
	0x0a, 0x09, 0x5f, 0x6d, 0x6e, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x42, 0x1c, 0x0a, 0x1a, 0x5f,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x6d, 0x69,
	0x6e, 0x47, 0x61, 0x73, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x42, 0x14, 0x0a,
	0x12, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x47, 0x61, 0x73, 0x41, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x22, 0xe4, 0x02, 0x0a, 0x0c, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f,
	0x6e, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f,
	0x6e, 0x69, 0x63, 0x12, 0x3c, 0x0a, 0x19, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x19, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c,
	0x64, 0x12, 0x28, 0x0a, 0x0f, 0x6d, 0x69, 0x6e, 0x47, 0x61, 0x73, 0x54, 0x68, 0x72, 0x65, 0x73,
	0x68, 0x6f, 0x6c, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x6d, 0x69, 0x6e, 0x47,
	0x61, 0x73, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x47, 0x61, 0x73, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x47, 0x61, 0x73, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x1a, 0x4f, 0x0a, 0x0b, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x6d, 0x0a, 0x06, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x2f, 0x0a, 0x06, 0x53, 0x65, 0x74, 0x52, 0x70, 0x63, 0x12, 0x0b,
	0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x07, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x70, 0x63,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x15, 0x5a, 0x13, 0x2e, 0x2f, 0x3b,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_src_proto_funds_service_config_proto_rawDescOnce sync.Once
	file_src_proto_funds_service_config_proto_rawDescData = file_src_proto_funds_service_config_proto_rawDesc
)

func file_src_proto_funds_service_config_proto_rawDescGZIP() []byte {
	file_src_proto_funds_service_config_proto_rawDescOnce.Do(func() {
		file_src_proto_funds_service_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_proto_funds_service_config_proto_rawDescData)
	})
	return file_src_proto_funds_service_config_proto_rawDescData
}

var file_src_proto_funds_service_config_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_src_proto_funds_service_config_proto_goTypes = []interface{}{
	(*Chain)(nil),         // 0: Chain
	(*SetRequest)(nil),    // 1: SetRequest
	(*LoadResponse)(nil),  // 2: LoadResponse
	(*Chain_Tron)(nil),    // 3: Chain.Tron
	nil,                   // 4: SetRequest.ChainsEntry
	nil,                   // 5: LoadResponse.ChainsEntry
	(*anypb.Any)(nil),     // 6: google.protobuf.Any
	(*emptypb.Empty)(nil), // 7: google.protobuf.Empty
}
var file_src_proto_funds_service_config_proto_depIdxs = []int32{
	4, // 0: SetRequest.chains:type_name -> SetRequest.ChainsEntry
	5, // 1: LoadResponse.chains:type_name -> LoadResponse.ChainsEntry
	6, // 2: SetRequest.ChainsEntry.value:type_name -> google.protobuf.Any
	6, // 3: LoadResponse.ChainsEntry.value:type_name -> google.protobuf.Any
	1, // 4: Config.SetRpc:input_type -> SetRequest
	7, // 5: Config.LoadRpc:input_type -> google.protobuf.Empty
	7, // 6: Config.SetRpc:output_type -> google.protobuf.Empty
	2, // 7: Config.LoadRpc:output_type -> LoadResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_src_proto_funds_service_config_proto_init() }
func file_src_proto_funds_service_config_proto_init() {
	if File_src_proto_funds_service_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_src_proto_funds_service_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chain); i {
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
		file_src_proto_funds_service_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetRequest); i {
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
		file_src_proto_funds_service_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadResponse); i {
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
		file_src_proto_funds_service_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chain_Tron); i {
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
	file_src_proto_funds_service_config_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_src_proto_funds_service_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_src_proto_funds_service_config_proto_goTypes,
		DependencyIndexes: file_src_proto_funds_service_config_proto_depIdxs,
		MessageInfos:      file_src_proto_funds_service_config_proto_msgTypes,
	}.Build()
	File_src_proto_funds_service_config_proto = out.File
	file_src_proto_funds_service_config_proto_rawDesc = nil
	file_src_proto_funds_service_config_proto_goTypes = nil
	file_src_proto_funds_service_config_proto_depIdxs = nil
}
