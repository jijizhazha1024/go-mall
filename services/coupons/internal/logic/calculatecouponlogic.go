package logic

import (
	"context"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateCouponLogic {
	return &CalculateCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CalculateCoupon 计算优惠卷折扣价格
func (l *CalculateCouponLogic) CalculateCoupon(in *coupons.CalculateCouponReq) (*coupons.CalculateCouponResp, error) {

	res := &coupons.CalculateCouponResp{}
	// 用户是否有改优惠券
	return res, nil
}
