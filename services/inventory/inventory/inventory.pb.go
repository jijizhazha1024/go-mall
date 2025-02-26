// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v3.19.4
// source: inventory.proto

package inventory

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

type InventoryReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*InventoryReq_Items  `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	PreOrderId    string                 `protobuf:"bytes,3,opt,name=pre_order_id,json=preOrderId,proto3" json:"pre_order_id,omitempty"`
	UserId        int32                  `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InventoryReq) Reset() {
	*x = InventoryReq{}
	mi := &file_inventory_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InventoryReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InventoryReq) ProtoMessage() {}

func (x *InventoryReq) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InventoryReq.ProtoReflect.Descriptor instead.
func (*InventoryReq) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{0}
}

func (x *InventoryReq) GetItems() []*InventoryReq_Items {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *InventoryReq) GetPreOrderId() string {
	if x != nil {
		return x.PreOrderId
	}
	return ""
}

func (x *InventoryReq) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetInventoryReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     int32                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetInventoryReq) Reset() {
	*x = GetInventoryReq{}
	mi := &file_inventory_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInventoryReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInventoryReq) ProtoMessage() {}

func (x *GetInventoryReq) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInventoryReq.ProtoReflect.Descriptor instead.
func (*GetInventoryReq) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{1}
}

func (x *GetInventoryReq) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

type BatchGetInventoryReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductIds    []int32                `protobuf:"varint,1,rep,packed,name=product_ids,json=productIds,proto3" json:"product_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BatchGetInventoryReq) Reset() {
	*x = BatchGetInventoryReq{}
	mi := &file_inventory_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchGetInventoryReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchGetInventoryReq) ProtoMessage() {}

func (x *BatchGetInventoryReq) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchGetInventoryReq.ProtoReflect.Descriptor instead.
func (*BatchGetInventoryReq) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{2}
}

func (x *BatchGetInventoryReq) GetProductIds() []int32 {
	if x != nil {
		return x.ProductIds
	}
	return nil
}

type PreInventoryRecord struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     int32                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity      int32                  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	ExpireTime    int64                  `protobuf:"varint,3,opt,name=expire_time,json=expireTime,proto3" json:"expire_time,omitempty"` // 过期时间戳
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PreInventoryRecord) Reset() {
	*x = PreInventoryRecord{}
	mi := &file_inventory_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PreInventoryRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreInventoryRecord) ProtoMessage() {}

func (x *PreInventoryRecord) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreInventoryRecord.ProtoReflect.Descriptor instead.
func (*PreInventoryRecord) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{3}
}

func (x *PreInventoryRecord) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *PreInventoryRecord) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *PreInventoryRecord) GetExpireTime() int64 {
	if x != nil {
		return x.ExpireTime
	}
	return 0
}

type InventoryResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StatusCode    int32                  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg     string                 `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InventoryResp) Reset() {
	*x = InventoryResp{}
	mi := &file_inventory_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InventoryResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InventoryResp) ProtoMessage() {}

func (x *InventoryResp) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InventoryResp.ProtoReflect.Descriptor instead.
func (*InventoryResp) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{4}
}

func (x *InventoryResp) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *InventoryResp) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

type GetInventoryResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StatusCode    int32                  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg     string                 `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	Inventory     int64                  `protobuf:"varint,3,opt,name=inventory,proto3" json:"inventory,omitempty"` // 返回操作之后库存数量
	SoldCount     int64                  `protobuf:"varint,4,opt,name=sold_count,json=soldCount,proto3" json:"sold_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetInventoryResp) Reset() {
	*x = GetInventoryResp{}
	mi := &file_inventory_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInventoryResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInventoryResp) ProtoMessage() {}

func (x *GetInventoryResp) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInventoryResp.ProtoReflect.Descriptor instead.
func (*GetInventoryResp) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{5}
}

func (x *GetInventoryResp) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *GetInventoryResp) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

func (x *GetInventoryResp) GetInventory() int64 {
	if x != nil {
		return x.Inventory
	}
	return 0
}

func (x *GetInventoryResp) GetSoldCount() int64 {
	if x != nil {
		return x.SoldCount
	}
	return 0
}

type InventoryReq_Items struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     int32                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity      int32                  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InventoryReq_Items) Reset() {
	*x = InventoryReq_Items{}
	mi := &file_inventory_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InventoryReq_Items) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InventoryReq_Items) ProtoMessage() {}

func (x *InventoryReq_Items) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InventoryReq_Items.ProtoReflect.Descriptor instead.
func (*InventoryReq_Items) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{0, 0}
}

func (x *InventoryReq_Items) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *InventoryReq_Items) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

var File_inventory_proto protoreflect.FileDescriptor

var file_inventory_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x22, 0xc2, 0x01, 0x0a,
	0x0c, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x12, 0x33, 0x0a,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x52, 0x65, 0x71, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x12, 0x20, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x65, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x1a, 0x42, 0x0a,
	0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x22, 0x30, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x49, 0x64, 0x22, 0x37, 0x0a, 0x14, 0x42, 0x61, 0x74, 0x63, 0x68, 0x47, 0x65, 0x74, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05,
	0x52, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x73, 0x22, 0x70, 0x0a, 0x12,
	0x50, 0x72, 0x65, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x1f, 0x0a,
	0x0b, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x4f,
	0x0a, 0x0d, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12,
	0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x73, 0x67, 0x22,
	0x8f, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f,
	0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x4d, 0x73, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x6f, 0x6c, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x6f, 0x6c, 0x64, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x32, 0xbc, 0x03, 0x0a, 0x09, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x47, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x1a, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x47, 0x65, 0x74, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x1b, 0x2e, 0x69, 0x6e,
	0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x44, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x17, 0x2e, 0x69, 0x6e,
	0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x49,
	0x0a, 0x14, 0x44, 0x65, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x50, 0x72, 0x65, 0x49, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x17, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x1a,
	0x18, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65,
	0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x46, 0x0a, 0x11, 0x44, 0x65, 0x63,
	0x72, 0x65, 0x61, 0x73, 0x65, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x17,
	0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x47, 0x0a, 0x12, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x50, 0x72, 0x65, 0x49, 0x6e,
	0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x17, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71,
	0x1a, 0x18, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x44, 0x0a, 0x0f, 0x52, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x17, 0x2e,
	0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74,
	0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inventory_proto_rawDescOnce sync.Once
	file_inventory_proto_rawDescData = file_inventory_proto_rawDesc
)

func file_inventory_proto_rawDescGZIP() []byte {
	file_inventory_proto_rawDescOnce.Do(func() {
		file_inventory_proto_rawDescData = protoimpl.X.CompressGZIP(file_inventory_proto_rawDescData)
	})
	return file_inventory_proto_rawDescData
}

var file_inventory_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_inventory_proto_goTypes = []any{
	(*InventoryReq)(nil),         // 0: inventory.InventoryReq
	(*GetInventoryReq)(nil),      // 1: inventory.GetInventoryReq
	(*BatchGetInventoryReq)(nil), // 2: inventory.BatchGetInventoryReq
	(*PreInventoryRecord)(nil),   // 3: inventory.PreInventoryRecord
	(*InventoryResp)(nil),        // 4: inventory.InventoryResp
	(*GetInventoryResp)(nil),     // 5: inventory.GetInventoryResp
	(*InventoryReq_Items)(nil),   // 6: inventory.InventoryReq.Items
}
var file_inventory_proto_depIdxs = []int32{
	6, // 0: inventory.InventoryReq.items:type_name -> inventory.InventoryReq.Items
	1, // 1: inventory.Inventory.GetInventory:input_type -> inventory.GetInventoryReq
	0, // 2: inventory.Inventory.UpdateInventory:input_type -> inventory.InventoryReq
	0, // 3: inventory.Inventory.DecreasePreInventory:input_type -> inventory.InventoryReq
	0, // 4: inventory.Inventory.DecreaseInventory:input_type -> inventory.InventoryReq
	0, // 5: inventory.Inventory.ReturnPreInventory:input_type -> inventory.InventoryReq
	0, // 6: inventory.Inventory.ReturnInventory:input_type -> inventory.InventoryReq
	5, // 7: inventory.Inventory.GetInventory:output_type -> inventory.GetInventoryResp
	4, // 8: inventory.Inventory.UpdateInventory:output_type -> inventory.InventoryResp
	4, // 9: inventory.Inventory.DecreasePreInventory:output_type -> inventory.InventoryResp
	4, // 10: inventory.Inventory.DecreaseInventory:output_type -> inventory.InventoryResp
	4, // 11: inventory.Inventory.ReturnPreInventory:output_type -> inventory.InventoryResp
	4, // 12: inventory.Inventory.ReturnInventory:output_type -> inventory.InventoryResp
	7, // [7:13] is the sub-list for method output_type
	1, // [1:7] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_inventory_proto_init() }
func file_inventory_proto_init() {
	if File_inventory_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_inventory_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_inventory_proto_goTypes,
		DependencyIndexes: file_inventory_proto_depIdxs,
		MessageInfos:      file_inventory_proto_msgTypes,
	}.Build()
	File_inventory_proto = out.File
	file_inventory_proto_rawDesc = nil
	file_inventory_proto_goTypes = nil
	file_inventory_proto_depIdxs = nil
}
