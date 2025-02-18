package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/biz"

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

	// param check
	if in.Pagination.Size <= 0 || in.Pagination.Page > biz.MaxPageSize {
		in.Pagination.Size = biz.MaxPageSize
	}
	if in.Pagination.Page <= 0 {
		in.Pagination.Page = 1
	}
	res := &coupons.ListUserCouponsResp{}
	userCoupons, err := l.svcCtx.UserCouponsModel.QueryUserCoupons(l.ctx, in.UserId, in.Pagination.Page, in.Pagination.Size)
	if err != nil {
		logx.Errorw("query user coupons error", logx.Field("err", err))
		return nil, err
	}
	res.UserCoupons = make([]*coupons.UserCoupon, 0, len(userCoupons))
	for _, uc := range userCoupons {
		res.UserCoupons = append(res.UserCoupons, convertUserCoupon2Resp(uc))
	}
	return res, nil
}
