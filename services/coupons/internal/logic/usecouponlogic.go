package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UseCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseCouponLogic {
	return &UseCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UseCoupon 使用优惠券（支付成功确认）
func (l *UseCouponLogic) UseCoupon(in *coupons.UseCouponReq) (*coupons.EmptyResp, error) {
	// 修改用户优惠券状态，记录优惠券使用记录，删除优惠券缓存
	return &coupons.EmptyResp{}, nil
}
