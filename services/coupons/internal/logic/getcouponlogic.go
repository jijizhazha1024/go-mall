package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCouponLogic {
	return &GetCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCoupon 获取单个优惠券
func (l *GetCouponLogic) GetCoupon(in *coupons.GetCouponReq) (*coupons.GetCouponResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.GetCouponResp{}, nil
}
