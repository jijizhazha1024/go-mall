package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClaimCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCouponLogic {
	return &ClaimCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClaimCoupon 用户领取优惠券
func (l *ClaimCouponLogic) ClaimCoupon(in *coupons.ClaimCouponReq) (*coupons.ClaimCouponResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.ClaimCouponResp{}, nil
}
