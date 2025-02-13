package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCouponUsagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCouponUsagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCouponUsagesLogic {
	return &ListCouponUsagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListCouponUsages 获取优惠券使用记录
func (l *ListCouponUsagesLogic) ListCouponUsages(in *coupons.ListCouponUsagesReq) (*coupons.ListCouponUsagesResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.ListCouponUsagesResp{}, nil
}
