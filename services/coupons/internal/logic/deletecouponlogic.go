package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCouponLogic {
	return &DeleteCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteCoupon 删除优惠券
func (l *DeleteCouponLogic) DeleteCoupon(in *coupons.DeleteCouponReq) (*coupons.DeleteCouponResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.DeleteCouponResp{}, nil
}
