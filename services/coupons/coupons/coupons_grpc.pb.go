// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.19.4
// source: coupons.proto

package coupons

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Coupons_ListCoupons_FullMethodName      = "/coupons.Coupons/ListCoupons"
	Coupons_GetCoupon_FullMethodName        = "/coupons.Coupons/GetCoupon"
	Coupons_CreateCoupon_FullMethodName     = "/coupons.Coupons/CreateCoupon"
	Coupons_UpdateCoupon_FullMethodName     = "/coupons.Coupons/UpdateCoupon"
	Coupons_DeleteCoupon_FullMethodName     = "/coupons.Coupons/DeleteCoupon"
	Coupons_ListUserCoupons_FullMethodName  = "/coupons.Coupons/ListUserCoupons"
	Coupons_GetUserCoupon_FullMethodName    = "/coupons.Coupons/GetUserCoupon"
	Coupons_ClaimCoupon_FullMethodName      = "/coupons.Coupons/ClaimCoupon"
	Coupons_UseCoupon_FullMethodName        = "/coupons.Coupons/UseCoupon"
	Coupons_ListCouponUsages_FullMethodName = "/coupons.Coupons/ListCouponUsages"
)

// CouponsClient is the client API for Coupons service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 优惠券服务
type CouponsClient interface {
	// ListCoupons 获取优惠券列表
	ListCoupons(ctx context.Context, in *ListCouponsReq, opts ...grpc.CallOption) (*ListCouponsResp, error)
	// GetCoupon 获取单个优惠券
	GetCoupon(ctx context.Context, in *GetCouponReq, opts ...grpc.CallOption) (*GetCouponResp, error)
	// CreateCoupon 创建优惠券
	CreateCoupon(ctx context.Context, in *CreateCouponReq, opts ...grpc.CallOption) (*CreateCouponResp, error)
	// UpdateCoupon 更新优惠券
	UpdateCoupon(ctx context.Context, in *UpdateCouponReq, opts ...grpc.CallOption) (*UpdateCouponResp, error)
	// DeleteCoupon 删除优惠券
	DeleteCoupon(ctx context.Context, in *DeleteCouponReq, opts ...grpc.CallOption) (*DeleteCouponResp, error)
	// ListUserCoupons 获取用户优惠券列表
	ListUserCoupons(ctx context.Context, in *ListUserCouponsReq, opts ...grpc.CallOption) (*ListUserCouponsResp, error)
	// GetUserCoupon 获取用户优惠券详情
	GetUserCoupon(ctx context.Context, in *GetUserCouponReq, opts ...grpc.CallOption) (*GetUserCouponResp, error)
	// ClaimCoupon 用户领取优惠券
	ClaimCoupon(ctx context.Context, in *ClaimCouponReq, opts ...grpc.CallOption) (*ClaimCouponResp, error)
	// UseCoupon 使用优惠券
	UseCoupon(ctx context.Context, in *UseCouponReq, opts ...grpc.CallOption) (*UseCouponResp, error)
	// ListCouponUsages 获取优惠券使用记录
	ListCouponUsages(ctx context.Context, in *ListCouponUsagesReq, opts ...grpc.CallOption) (*ListCouponUsagesResp, error)
}

type couponsClient struct {
	cc grpc.ClientConnInterface
}

func NewCouponsClient(cc grpc.ClientConnInterface) CouponsClient {
	return &couponsClient{cc}
}

func (c *couponsClient) ListCoupons(ctx context.Context, in *ListCouponsReq, opts ...grpc.CallOption) (*ListCouponsResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCouponsResp)
	err := c.cc.Invoke(ctx, Coupons_ListCoupons_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) GetCoupon(ctx context.Context, in *GetCouponReq, opts ...grpc.CallOption) (*GetCouponResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCouponResp)
	err := c.cc.Invoke(ctx, Coupons_GetCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) CreateCoupon(ctx context.Context, in *CreateCouponReq, opts ...grpc.CallOption) (*CreateCouponResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCouponResp)
	err := c.cc.Invoke(ctx, Coupons_CreateCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) UpdateCoupon(ctx context.Context, in *UpdateCouponReq, opts ...grpc.CallOption) (*UpdateCouponResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateCouponResp)
	err := c.cc.Invoke(ctx, Coupons_UpdateCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) DeleteCoupon(ctx context.Context, in *DeleteCouponReq, opts ...grpc.CallOption) (*DeleteCouponResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCouponResp)
	err := c.cc.Invoke(ctx, Coupons_DeleteCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) ListUserCoupons(ctx context.Context, in *ListUserCouponsReq, opts ...grpc.CallOption) (*ListUserCouponsResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListUserCouponsResp)
	err := c.cc.Invoke(ctx, Coupons_ListUserCoupons_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) GetUserCoupon(ctx context.Context, in *GetUserCouponReq, opts ...grpc.CallOption) (*GetUserCouponResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserCouponResp)
	err := c.cc.Invoke(ctx, Coupons_GetUserCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) ClaimCoupon(ctx context.Context, in *ClaimCouponReq, opts ...grpc.CallOption) (*ClaimCouponResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClaimCouponResp)
	err := c.cc.Invoke(ctx, Coupons_ClaimCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) UseCoupon(ctx context.Context, in *UseCouponReq, opts ...grpc.CallOption) (*UseCouponResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UseCouponResp)
	err := c.cc.Invoke(ctx, Coupons_UseCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponsClient) ListCouponUsages(ctx context.Context, in *ListCouponUsagesReq, opts ...grpc.CallOption) (*ListCouponUsagesResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCouponUsagesResp)
	err := c.cc.Invoke(ctx, Coupons_ListCouponUsages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CouponsServer is the server API for Coupons service.
// All implementations must embed UnimplementedCouponsServer
// for forward compatibility
//
// 优惠券服务
type CouponsServer interface {
	// ListCoupons 获取优惠券列表
	ListCoupons(context.Context, *ListCouponsReq) (*ListCouponsResp, error)
	// GetCoupon 获取单个优惠券
	GetCoupon(context.Context, *GetCouponReq) (*GetCouponResp, error)
	// CreateCoupon 创建优惠券
	CreateCoupon(context.Context, *CreateCouponReq) (*CreateCouponResp, error)
	// UpdateCoupon 更新优惠券
	UpdateCoupon(context.Context, *UpdateCouponReq) (*UpdateCouponResp, error)
	// DeleteCoupon 删除优惠券
	DeleteCoupon(context.Context, *DeleteCouponReq) (*DeleteCouponResp, error)
	// ListUserCoupons 获取用户优惠券列表
	ListUserCoupons(context.Context, *ListUserCouponsReq) (*ListUserCouponsResp, error)
	// GetUserCoupon 获取用户优惠券详情
	GetUserCoupon(context.Context, *GetUserCouponReq) (*GetUserCouponResp, error)
	// ClaimCoupon 用户领取优惠券
	ClaimCoupon(context.Context, *ClaimCouponReq) (*ClaimCouponResp, error)
	// UseCoupon 使用优惠券
	UseCoupon(context.Context, *UseCouponReq) (*UseCouponResp, error)
	// ListCouponUsages 获取优惠券使用记录
	ListCouponUsages(context.Context, *ListCouponUsagesReq) (*ListCouponUsagesResp, error)
	mustEmbedUnimplementedCouponsServer()
}

// UnimplementedCouponsServer must be embedded to have forward compatible implementations.
type UnimplementedCouponsServer struct {
}

func (UnimplementedCouponsServer) ListCoupons(context.Context, *ListCouponsReq) (*ListCouponsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCoupons not implemented")
}
func (UnimplementedCouponsServer) GetCoupon(context.Context, *GetCouponReq) (*GetCouponResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCoupon not implemented")
}
func (UnimplementedCouponsServer) CreateCoupon(context.Context, *CreateCouponReq) (*CreateCouponResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCoupon not implemented")
}
func (UnimplementedCouponsServer) UpdateCoupon(context.Context, *UpdateCouponReq) (*UpdateCouponResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCoupon not implemented")
}
func (UnimplementedCouponsServer) DeleteCoupon(context.Context, *DeleteCouponReq) (*DeleteCouponResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCoupon not implemented")
}
func (UnimplementedCouponsServer) ListUserCoupons(context.Context, *ListUserCouponsReq) (*ListUserCouponsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserCoupons not implemented")
}
func (UnimplementedCouponsServer) GetUserCoupon(context.Context, *GetUserCouponReq) (*GetUserCouponResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserCoupon not implemented")
}
func (UnimplementedCouponsServer) ClaimCoupon(context.Context, *ClaimCouponReq) (*ClaimCouponResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimCoupon not implemented")
}
func (UnimplementedCouponsServer) UseCoupon(context.Context, *UseCouponReq) (*UseCouponResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseCoupon not implemented")
}
func (UnimplementedCouponsServer) ListCouponUsages(context.Context, *ListCouponUsagesReq) (*ListCouponUsagesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCouponUsages not implemented")
}
func (UnimplementedCouponsServer) mustEmbedUnimplementedCouponsServer() {}

// UnsafeCouponsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CouponsServer will
// result in compilation errors.
type UnsafeCouponsServer interface {
	mustEmbedUnimplementedCouponsServer()
}

func RegisterCouponsServer(s grpc.ServiceRegistrar, srv CouponsServer) {
	s.RegisterService(&Coupons_ServiceDesc, srv)
}

func _Coupons_ListCoupons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCouponsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).ListCoupons(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_ListCoupons_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).ListCoupons(ctx, req.(*ListCouponsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_GetCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCouponReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).GetCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_GetCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).GetCoupon(ctx, req.(*GetCouponReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_CreateCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCouponReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).CreateCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_CreateCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).CreateCoupon(ctx, req.(*CreateCouponReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_UpdateCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCouponReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).UpdateCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_UpdateCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).UpdateCoupon(ctx, req.(*UpdateCouponReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_DeleteCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCouponReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).DeleteCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_DeleteCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).DeleteCoupon(ctx, req.(*DeleteCouponReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_ListUserCoupons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserCouponsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).ListUserCoupons(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_ListUserCoupons_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).ListUserCoupons(ctx, req.(*ListUserCouponsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_GetUserCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserCouponReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).GetUserCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_GetUserCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).GetUserCoupon(ctx, req.(*GetUserCouponReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_ClaimCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClaimCouponReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).ClaimCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_ClaimCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).ClaimCoupon(ctx, req.(*ClaimCouponReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_UseCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseCouponReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).UseCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_UseCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).UseCoupon(ctx, req.(*UseCouponReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coupons_ListCouponUsages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCouponUsagesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).ListCouponUsages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coupons_ListCouponUsages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).ListCouponUsages(ctx, req.(*ListCouponUsagesReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Coupons_ServiceDesc is the grpc.ServiceDesc for Coupons service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Coupons_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "coupons.Coupons",
	HandlerType: (*CouponsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCoupons",
			Handler:    _Coupons_ListCoupons_Handler,
		},
		{
			MethodName: "GetCoupon",
			Handler:    _Coupons_GetCoupon_Handler,
		},
		{
			MethodName: "CreateCoupon",
			Handler:    _Coupons_CreateCoupon_Handler,
		},
		{
			MethodName: "UpdateCoupon",
			Handler:    _Coupons_UpdateCoupon_Handler,
		},
		{
			MethodName: "DeleteCoupon",
			Handler:    _Coupons_DeleteCoupon_Handler,
		},
		{
			MethodName: "ListUserCoupons",
			Handler:    _Coupons_ListUserCoupons_Handler,
		},
		{
			MethodName: "GetUserCoupon",
			Handler:    _Coupons_GetUserCoupon_Handler,
		},
		{
			MethodName: "ClaimCoupon",
			Handler:    _Coupons_ClaimCoupon_Handler,
		},
		{
			MethodName: "UseCoupon",
			Handler:    _Coupons_UseCoupon_Handler,
		},
		{
			MethodName: "ListCouponUsages",
			Handler:    _Coupons_ListCouponUsages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coupons.proto",
}
