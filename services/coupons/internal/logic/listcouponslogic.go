package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCouponsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCouponsLogic {
	return &ListCouponsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListCoupons 获取优惠券列表
func (l *ListCouponsLogic) ListCoupons(in *coupons.ListCouponsReq) (*coupons.ListCouponsResp, error) {
	// todo: add your logic here and delete this line

	return &coupons.ListCouponsResp{}, nil
}
