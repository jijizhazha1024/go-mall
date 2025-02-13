package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserCouponsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserCouponsLogic {
	return &ListUserCouponsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListUserCoupons 获取用户优惠券列表
func (l *ListUserCouponsLogic) ListUserCoupons(in *coupons.ListUserCouponsReq) (*coupons.ListUserCouponsResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.ListUserCouponsResp{}, nil
}
