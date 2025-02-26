package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"
	"time"

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
	// param check
	if in.Pagination.Size <= 0 || in.Pagination.Size > biz.MaxPageSize {
		in.Pagination.Page = biz.MaxPageSize
	}
	if in.Pagination.Page <= 0 {
		in.Pagination.Page = 1
	}

	couponsUsageList, err := l.svcCtx.CouponUsageModel.QueryUsageListByUserId(l.ctx, uint64(in.UserId), in.Pagination.Page, in.Pagination.Size)
	if err != nil {
		logx.Errorw("query coupon usage error", logx.Field("err", err))
		return nil, err
	}
	res := &coupons.ListCouponUsagesResp{
		Usages: make([]*coupons.CouponUsage, 0, len(couponsUsageList)),
	}

	for _, couponUsage := range couponsUsageList {
		res.Usages = append(res.Usages, &coupons.CouponUsage{
			Id:         int32(couponUsage.Id),
			CouponId:   couponUsage.CouponId,
			CouponType: coupons.CouponType(couponUsage.CouponType),
			// 确保浮点数精度
			OriginValue:    convertToYuan(couponUsage.OriginValue),
			DiscountAmount: convertToYuan(couponUsage.DiscountAmount),
			OrderId:        couponUsage.OrderId,
			UserId:         int32(couponUsage.UserId),
			AppliedAt:      couponUsage.AppliedAt.Format(time.DateTime),
		})
	}
	res.TotalCount = int32(len(couponsUsageList))

	return res, nil
}
