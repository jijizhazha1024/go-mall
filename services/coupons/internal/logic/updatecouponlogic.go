package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCouponLogic {
	return &UpdateCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateCoupon 更新优惠券
func (l *UpdateCouponLogic) UpdateCoupon(in *coupons.UpdateCouponReq) (*coupons.UpdateCouponResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.UpdateCouponResp{}, nil
}
