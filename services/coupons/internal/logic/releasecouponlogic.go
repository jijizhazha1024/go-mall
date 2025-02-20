package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseCouponLogic {
	return &ReleaseCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReleaseCoupon 释放优惠券（订单取消/超时释放）
func (l *ReleaseCouponLogic) ReleaseCoupon(in *coupons.ReleaseCouponReq) (*coupons.EmptyResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.EmptyResp{}, nil
}
