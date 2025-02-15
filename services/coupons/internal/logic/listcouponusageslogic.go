package logic

import (
	"context"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"
	"strconv"
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
	couponsUsageList, err := l.svcCtx.CouponUsageModel.QueryUsageListByUserId(l.ctx, uint64(in.UserId), in.Pagination.Page, in.Pagination.Limit)
	if err != nil {
		logx.Errorw("query coupon usage error", logx.Field("err", err))
		return nil, err
	}
	res := &coupons.ListCouponUsagesResp{
		Usages: make([]*coupons.CouponUsage, 0, len(couponsUsageList)),
	}
	//
	for _, couponUsage := range couponsUsageList {
		res.Usages = append(res.Usages, &coupons.CouponUsage{
			Id:         int32(couponUsage.Id),
			CouponId:   couponUsage.CouponId,
			CouponType: coupons.CouponType(couponUsage.CouponType),
			// 确保浮点数精度
			OriginValue:    strconv.FormatFloat(couponUsage.OriginValue, 'f', 2, 64),
			DiscountAmount: strconv.FormatFloat(couponUsage.DiscountAmount, 'f', 2, 64),
			OrderId:        couponUsage.OrderId,
			UserId:         int32(couponUsage.UserId),
			AppliedAt:      couponUsage.AppliedAt.Format(time.DateTime),
		})
	}
	res.TotalCount = int32(len(couponsUsageList))

	return &coupons.ListCouponUsagesResp{}, nil
}
