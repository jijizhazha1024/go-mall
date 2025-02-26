// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: checkout.proto

package checkout

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	CheckoutService_PrepareCheckout_FullMethodName            = "/checkout.CheckoutService/PrepareCheckout"
	CheckoutService_ReleaseCheckout_FullMethodName            = "/checkout.CheckoutService/ReleaseCheckout"
	CheckoutService_GetCheckoutList_FullMethodName            = "/checkout.CheckoutService/GetCheckoutList"
	CheckoutService_GetCheckoutDetail_FullMethodName          = "/checkout.CheckoutService/GetCheckoutDetail"
	CheckoutService_UpdateStatus2Order_FullMethodName         = "/checkout.CheckoutService/UpdateStatus2Order"
	CheckoutService_UpdateStatus2OrderRollback_FullMethodName = "/checkout.CheckoutService/UpdateStatus2OrderRollback"
)

// CheckoutServiceClient is the client API for CheckoutService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckoutServiceClient interface {
	// PrepareCheckout 预结算)生成预订单）
	PrepareCheckout(ctx context.Context, in *CheckoutReq, opts ...grpc.CallOption) (*CheckoutResp, error)
	// UpdateCheckoutStatus2Success 当订单超时，支付超时，支付退款
	ReleaseCheckout(ctx context.Context, in *ReleaseReq, opts ...grpc.CallOption) (*EmptyResp, error)
	// GetCheckoutList 获取结算列表
	GetCheckoutList(ctx context.Context, in *CheckoutListReq, opts ...grpc.CallOption) (*CheckoutListResp, error)
	// GetCheckoutDetail 获取结算详情
	GetCheckoutDetail(ctx context.Context, in *CheckoutDetailReq, opts ...grpc.CallOption) (*CheckoutDetailResp, error)
	// UpdateStatus2Order 由订单服务调用，更新结算状态为已确认
	UpdateStatus2Order(ctx context.Context, in *UpdateStatusReq, opts ...grpc.CallOption) (*EmptyResp, error)
	// UpdateStatus2OrderRollback 补偿操作
	UpdateStatus2OrderRollback(ctx context.Context, in *UpdateStatusReq, opts ...grpc.CallOption) (*EmptyResp, error)
}

type checkoutServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCheckoutServiceClient(cc grpc.ClientConnInterface) CheckoutServiceClient {
	return &checkoutServiceClient{cc}
}

func (c *checkoutServiceClient) PrepareCheckout(ctx context.Context, in *CheckoutReq, opts ...grpc.CallOption) (*CheckoutResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckoutResp)
	err := c.cc.Invoke(ctx, CheckoutService_PrepareCheckout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) ReleaseCheckout(ctx context.Context, in *ReleaseReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResp)
	err := c.cc.Invoke(ctx, CheckoutService_ReleaseCheckout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) GetCheckoutList(ctx context.Context, in *CheckoutListReq, opts ...grpc.CallOption) (*CheckoutListResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckoutListResp)
	err := c.cc.Invoke(ctx, CheckoutService_GetCheckoutList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) GetCheckoutDetail(ctx context.Context, in *CheckoutDetailReq, opts ...grpc.CallOption) (*CheckoutDetailResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckoutDetailResp)
	err := c.cc.Invoke(ctx, CheckoutService_GetCheckoutDetail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) UpdateStatus2Order(ctx context.Context, in *UpdateStatusReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResp)
	err := c.cc.Invoke(ctx, CheckoutService_UpdateStatus2Order_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) UpdateStatus2OrderRollback(ctx context.Context, in *UpdateStatusReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResp)
	err := c.cc.Invoke(ctx, CheckoutService_UpdateStatus2OrderRollback_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckoutServiceServer is the server API for CheckoutService service.
// All implementations must embed UnimplementedCheckoutServiceServer
// for forward compatibility.
type CheckoutServiceServer interface {
	// PrepareCheckout 预结算)生成预订单）
	PrepareCheckout(context.Context, *CheckoutReq) (*CheckoutResp, error)
	// UpdateCheckoutStatus2Success 当订单超时，支付超时，支付退款
	ReleaseCheckout(context.Context, *ReleaseReq) (*EmptyResp, error)
	// GetCheckoutList 获取结算列表
	GetCheckoutList(context.Context, *CheckoutListReq) (*CheckoutListResp, error)
	// GetCheckoutDetail 获取结算详情
	GetCheckoutDetail(context.Context, *CheckoutDetailReq) (*CheckoutDetailResp, error)
	// UpdateStatus2Order 由订单服务调用，更新结算状态为已确认
	UpdateStatus2Order(context.Context, *UpdateStatusReq) (*EmptyResp, error)
	// UpdateStatus2OrderRollback 补偿操作
	UpdateStatus2OrderRollback(context.Context, *UpdateStatusReq) (*EmptyResp, error)
	mustEmbedUnimplementedCheckoutServiceServer()
}

// UnimplementedCheckoutServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCheckoutServiceServer struct{}

func (UnimplementedCheckoutServiceServer) PrepareCheckout(context.Context, *CheckoutReq) (*CheckoutResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrepareCheckout not implemented")
}
func (UnimplementedCheckoutServiceServer) ReleaseCheckout(context.Context, *ReleaseReq) (*EmptyResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseCheckout not implemented")
}
func (UnimplementedCheckoutServiceServer) GetCheckoutList(context.Context, *CheckoutListReq) (*CheckoutListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCheckoutList not implemented")
}
func (UnimplementedCheckoutServiceServer) GetCheckoutDetail(context.Context, *CheckoutDetailReq) (*CheckoutDetailResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCheckoutDetail not implemented")
}
func (UnimplementedCheckoutServiceServer) UpdateStatus2Order(context.Context, *UpdateStatusReq) (*EmptyResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus2Order not implemented")
}
func (UnimplementedCheckoutServiceServer) UpdateStatus2OrderRollback(context.Context, *UpdateStatusReq) (*EmptyResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus2OrderRollback not implemented")
}
func (UnimplementedCheckoutServiceServer) mustEmbedUnimplementedCheckoutServiceServer() {}
func (UnimplementedCheckoutServiceServer) testEmbeddedByValue()                         {}

// UnsafeCheckoutServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CheckoutServiceServer will
// result in compilation errors.
type UnsafeCheckoutServiceServer interface {
	mustEmbedUnimplementedCheckoutServiceServer()
}

func RegisterCheckoutServiceServer(s grpc.ServiceRegistrar, srv CheckoutServiceServer) {
	// If the following call pancis, it indicates UnimplementedCheckoutServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CheckoutService_ServiceDesc, srv)
}

func _CheckoutService_PrepareCheckout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckoutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).PrepareCheckout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckoutService_PrepareCheckout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).PrepareCheckout(ctx, req.(*CheckoutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_ReleaseCheckout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).ReleaseCheckout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckoutService_ReleaseCheckout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).ReleaseCheckout(ctx, req.(*ReleaseReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_GetCheckoutList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckoutListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).GetCheckoutList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckoutService_GetCheckoutList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).GetCheckoutList(ctx, req.(*CheckoutListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_GetCheckoutDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckoutDetailReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).GetCheckoutDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckoutService_GetCheckoutDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).GetCheckoutDetail(ctx, req.(*CheckoutDetailReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_UpdateStatus2Order_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).UpdateStatus2Order(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckoutService_UpdateStatus2Order_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).UpdateStatus2Order(ctx, req.(*UpdateStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_UpdateStatus2OrderRollback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).UpdateStatus2OrderRollback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckoutService_UpdateStatus2OrderRollback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).UpdateStatus2OrderRollback(ctx, req.(*UpdateStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

// CheckoutService_ServiceDesc is the grpc.ServiceDesc for CheckoutService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CheckoutService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "checkout.CheckoutService",
	HandlerType: (*CheckoutServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PrepareCheckout",
			Handler:    _CheckoutService_PrepareCheckout_Handler,
		},
		{
			MethodName: "ReleaseCheckout",
			Handler:    _CheckoutService_ReleaseCheckout_Handler,
		},
		{
			MethodName: "GetCheckoutList",
			Handler:    _CheckoutService_GetCheckoutList_Handler,
		},
		{
			MethodName: "GetCheckoutDetail",
			Handler:    _CheckoutService_GetCheckoutDetail_Handler,
		},
		{
			MethodName: "UpdateStatus2Order",
			Handler:    _CheckoutService_UpdateStatus2Order_Handler,
		},
		{
			MethodName: "UpdateStatus2OrderRollback",
			Handler:    _CheckoutService_UpdateStatus2OrderRollback_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "checkout.proto",
}
