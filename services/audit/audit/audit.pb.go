// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.19.4
// source: audit.proto

package audit

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

// 创建审计日志请求
type CreateAuditLogReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId            uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`            // 执行操作的用户ID
	Username          string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`                       // 执行操作的用户名
	ActionType        string `protobuf:"bytes,3,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"` // 操作类型（新增、修改、删除等）
	ActionDescription string `protobuf:"bytes,4,opt,name=action_description,json=actionDescription,proto3" json:"action_description,omitempty"`
	TargetTable       string `protobuf:"bytes,5,opt,name=target_table,json=targetTable,proto3" json:"target_table,omitempty"` // 目标类型（如：商品、订单等）
	TargetId          int64  `protobuf:"varint,6,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`         // 目标对象的ID（如商品ID、订单ID等）
	OldData           string `protobuf:"bytes,7,opt,name=old_data,json=oldData,proto3" json:"old_data,omitempty"`             // 操作前的数据（JSON格式）
	NewData           string `protobuf:"bytes,8,opt,name=new_data,json=newData,proto3" json:"new_data,omitempty"`             // 操作后的数据（JSON格式）
	CreateAt          int64  `protobuf:"varint,9,opt,name=create_at,json=createAt,proto3" json:"create_at,omitempty"`
}

func (x *CreateAuditLogReq) Reset() {
	*x = CreateAuditLogReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAuditLogReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAuditLogReq) ProtoMessage() {}

func (x *CreateAuditLogReq) ProtoReflect() protoreflect.Message {
	mi := &file_audit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAuditLogReq.ProtoReflect.Descriptor instead.
func (*CreateAuditLogReq) Descriptor() ([]byte, []int) {
	return file_audit_proto_rawDescGZIP(), []int{0}
}

func (x *CreateAuditLogReq) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateAuditLogReq) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *CreateAuditLogReq) GetActionType() string {
	if x != nil {
		return x.ActionType
	}
	return ""
}

func (x *CreateAuditLogReq) GetActionDescription() string {
	if x != nil {
		return x.ActionDescription
	}
	return ""
}

func (x *CreateAuditLogReq) GetTargetTable() string {
	if x != nil {
		return x.TargetTable
	}
	return ""
}

func (x *CreateAuditLogReq) GetTargetId() int64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *CreateAuditLogReq) GetOldData() string {
	if x != nil {
		return x.OldData
	}
	return ""
}

func (x *CreateAuditLogReq) GetNewData() string {
	if x != nil {
		return x.NewData
	}
	return ""
}

func (x *CreateAuditLogReq) GetCreateAt() int64 {
	if x != nil {
		return x.CreateAt
	}
	return 0
}

// 创建审计日志响应
type CreateAuditLogRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode uint32 `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	Ok         bool   `protobuf:"varint,3,opt,name=ok,proto3" json:"ok,omitempty"` // 审计日志的ID
}

func (x *CreateAuditLogRes) Reset() {
	*x = CreateAuditLogRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_audit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAuditLogRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAuditLogRes) ProtoMessage() {}

func (x *CreateAuditLogRes) ProtoReflect() protoreflect.Message {
	mi := &file_audit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAuditLogRes.ProtoReflect.Descriptor instead.
func (*CreateAuditLogRes) Descriptor() ([]byte, []int) {
	return file_audit_proto_rawDescGZIP(), []int{1}
}

func (x *CreateAuditLogRes) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *CreateAuditLogRes) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

func (x *CreateAuditLogRes) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_audit_proto protoreflect.FileDescriptor

var file_audit_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x22, 0xab, 0x02, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41,
	0x75, 0x64, 0x69, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x2d, 0x0a, 0x12, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x21, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x12,
	0x19, 0x0a, 0x08, 0x6f, 0x6c, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6f, 0x6c, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x65,
	0x77, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65,
	0x77, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f,
	0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x41, 0x74, 0x22, 0x63, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x64, 0x69,
	0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x4d, 0x73, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x32, 0x4d, 0x0a, 0x05, 0x41, 0x75, 0x64, 0x69, 0x74,
	0x12, 0x44, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x64, 0x69, 0x74, 0x4c,
	0x6f, 0x67, 0x12, 0x18, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x41, 0x75, 0x64, 0x69, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x64, 0x69, 0x74,
	0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x61, 0x75, 0x64, 0x69,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_audit_proto_rawDescOnce sync.Once
	file_audit_proto_rawDescData = file_audit_proto_rawDesc
)

func file_audit_proto_rawDescGZIP() []byte {
	file_audit_proto_rawDescOnce.Do(func() {
		file_audit_proto_rawDescData = protoimpl.X.CompressGZIP(file_audit_proto_rawDescData)
	})
	return file_audit_proto_rawDescData
}

var file_audit_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_audit_proto_goTypes = []any{
	(*CreateAuditLogReq)(nil), // 0: audit.CreateAuditLogReq
	(*CreateAuditLogRes)(nil), // 1: audit.CreateAuditLogRes
}
var file_audit_proto_depIdxs = []int32{
	0, // 0: audit.Audit.CreateAuditLog:input_type -> audit.CreateAuditLogReq
	1, // 1: audit.Audit.CreateAuditLog:output_type -> audit.CreateAuditLogRes
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_audit_proto_init() }
func file_audit_proto_init() {
	if File_audit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_audit_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateAuditLogReq); i {
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
		file_audit_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateAuditLogRes); i {
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
			RawDescriptor: file_audit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_audit_proto_goTypes,
		DependencyIndexes: file_audit_proto_depIdxs,
		MessageInfos:      file_audit_proto_msgTypes,
	}.Build()
	File_audit_proto = out.File
	file_audit_proto_rawDesc = nil
	file_audit_proto_goTypes = nil
	file_audit_proto_depIdxs = nil
}
