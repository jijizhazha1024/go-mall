package logic

import (
	"github.com/shopspring/decimal"
	"jijizhazha1024/go-mall/dal/model/coupons/coupon"
	"jijizhazha1024/go-mall/dal/model/coupons/user_coupons"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"time"
)

func convertCoupon2Resp(c *coupon.Coupons) *coupons.Coupon {
	return &coupons.Coupon{
		Id:             c.Id,
		Name:           c.Name,
		Type:           coupons.CouponType(c.Type),
		Value:          convertToYuan(c.Value),
		MinAmount:      convertToYuan(c.MinAmount),
		TotalCount:     c.TotalCount, // 发放
		RemainingCount: c.RemainingCount,
		StartTime:      c.StartTime.Format(time.DateTime),
		EndTime:        c.EndTime.Format(time.DateTime),
		CreatedAt:      c.CreatedAt.Format(time.DateTime),
		UpdatedAt:      c.UpdatedAt.Format(time.DateTime),
	}
}

func convertUserCoupon2Resp(uc *user_coupons.UserCoupons) *coupons.UserCoupon {
	return &coupons.UserCoupon{
		Id:        int32(uc.Id),
		UserId:    int32(uc.UserId),
		CouponId:  uc.CouponId,
		Status:    coupons.CouponStatus(uc.Status),
		OrderId:   uc.OrderId.String,
		UsedAt:    uc.UsedAt.Time.Format(time.DateTime),
		CreatedAt: uc.CreatedAt.Format(time.DateTime),
		UpdatedAt: uc.UpdatedAt.Format(time.DateTime),
	}
}

func convertToYuan(fen int64) string {
	return decimal.NewFromInt(fen).
		Div(decimal.NewFromInt(100)).
		StringFixedBank(2) // 银行家舍入法
}
