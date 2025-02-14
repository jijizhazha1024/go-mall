package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"jijizhazha1024/go-mall/common/consts/biz"

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
	// param check
	if in.Pagination.Limit <= 0 || in.Pagination.Limit > biz.MaxPageSize {
		in.Pagination.Limit = biz.MaxPageSize
	}
	if in.Pagination.Page <= 0 {
		in.Pagination.Page = 1
	}
	res := &coupons.ListCouponsResp{}
	queryCoupons, err := l.svcCtx.CouponsModel.QueryCoupons(l.ctx, in.Pagination.Page, in.Pagination.Limit, in.Type)
	if err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return res, nil
		}
		logx.Errorw("query coupons error", logx.Field("err", err))
		return nil, err
	}
	couponsList := make([]*coupons.Coupon, 0, len(queryCoupons))
	for _, c := range queryCoupons {
		couponsList = append(couponsList, convertCoupon2Resp(c))
	}
	res.Coupons = couponsList
	return res, nil
}
