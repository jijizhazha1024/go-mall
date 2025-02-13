package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCouponLogic {
	return &CreateCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateCoupon 创建优惠券
func (l *CreateCouponLogic) CreateCoupon(in *coupons.CreateCouponReq) (*coupons.CreateCouponResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.CreateCouponResp{}, nil
}
