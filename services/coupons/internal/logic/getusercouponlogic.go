package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCouponLogic {
	return &GetUserCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserCoupon 获取用户优惠券详情
func (l *GetUserCouponLogic) GetUserCoupon(in *coupons.GetUserCouponReq) (*coupons.GetUserCouponResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.GetUserCouponResp{}, nil
}
