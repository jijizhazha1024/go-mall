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
	for _, coupon := range queryCoupons {
		res.Coupons = append(res.Coupons, &coupons.Coupon{
			Id:             coupon.Id,
			Name:           coupon.Name,
			Type:           coupons.CouponType(coupon.Type),
			Value:          int32(coupon.Value),
			MinAmount:      int32(coupon.MinAmount),
			RemainingCount: int32(coupon.RemainingCount),
			TotalCount:     int32(coupon.TotalCount),
			EndTime:        coupon.EndTime.Format("2006-01-02 15:04:05"),
			StartTime:      coupon.StartTime.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}
