// Code generated by goctl. DO NOT EDIT.
// Source: coupons.proto

package server

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/logic"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"
)

type CouponsServer struct {
	svcCtx *svc.ServiceContext
	coupons.UnimplementedCouponsServer
}

func NewCouponsServer(svcCtx *svc.ServiceContext) *CouponsServer {
	return &CouponsServer{
		svcCtx: svcCtx,
	}
}

// ListCoupons 获取优惠券列表
func (s *CouponsServer) ListCoupons(ctx context.Context, in *coupons.ListCouponsReq) (*coupons.ListCouponsResp, error) {
	l := logic.NewListCouponsLogic(ctx, s.svcCtx)
	return l.ListCoupons(in)
}

// GetCoupon 获取单个优惠券
func (s *CouponsServer) GetCoupon(ctx context.Context, in *coupons.GetCouponReq) (*coupons.GetCouponResp, error) {
	l := logic.NewGetCouponLogic(ctx, s.svcCtx)
	return l.GetCoupon(in)
}

// ClaimCoupon 用户领取优惠券
func (s *CouponsServer) ClaimCoupon(ctx context.Context, in *coupons.ClaimCouponReq) (*coupons.ClaimCouponResp, error) {
	l := logic.NewClaimCouponLogic(ctx, s.svcCtx)
	return l.ClaimCoupon(in)
}

// ListUserCoupons 获取用户优惠券列表
func (s *CouponsServer) ListUserCoupons(ctx context.Context, in *coupons.ListUserCouponsReq) (*coupons.ListUserCouponsResp, error) {
	l := logic.NewListUserCouponsLogic(ctx, s.svcCtx)
	return l.ListUserCoupons(in)
}

// CalculateCoupon 计算优惠券
func (s *CouponsServer) CalculateCoupon(ctx context.Context, in *coupons.CalculateCouponReq) (*coupons.CalculateCouponResp, error) {
	l := logic.NewCalculateCouponLogic(ctx, s.svcCtx)
	return l.CalculateCoupon(in)
}

// ListCouponUsages 获取优惠券使用记录
func (s *CouponsServer) ListCouponUsages(ctx context.Context, in *coupons.ListCouponUsagesReq) (*coupons.ListCouponUsagesResp, error) {
	l := logic.NewListCouponUsagesLogic(ctx, s.svcCtx)
	return l.ListCouponUsages(in)
}

// --------------- 使用优惠券 --------------- pre_order_id来进行使用
func (s *CouponsServer) LockCoupon(ctx context.Context, in *coupons.LockCouponReq) (*coupons.EmptyResp, error) {
	l := logic.NewLockCouponLogic(ctx, s.svcCtx)
	return l.LockCoupon(in)
}

// 释放优惠券（订单取消/超时释放）
func (s *CouponsServer) ReleaseCoupon(ctx context.Context, in *coupons.ReleaseCouponReq) (*coupons.EmptyResp, error) {
	l := logic.NewReleaseCouponLogic(ctx, s.svcCtx)
	return l.ReleaseCoupon(in)
}

// 使用优惠券（支付成功确认）
func (s *CouponsServer) UseCoupon(ctx context.Context, in *coupons.UseCouponReq) (*coupons.EmptyResp, error) {
	l := logic.NewUseCouponLogic(ctx, s.svcCtx)
	return l.UseCoupon(in)
}
