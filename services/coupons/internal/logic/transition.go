package logic

import (
	"jijizhazha1024/go-mall/dal/model/coupons/coupon"
	"jijizhazha1024/go-mall/services/coupons/coupons"
)

func convertCoupon2Resp(c *coupon.Coupons) *coupons.Coupon {
	return &coupons.Coupon{
		Id:             c.Id,
		Name:           c.Name,
		Type:           coupons.CouponType(c.Type),
		Value:          int32(c.Value),
		MinAmount:      int32(c.MinAmount),
		TotalCount:     int32(c.TotalCount),
		RemainingCount: int32(c.RemainingCount),
		StartTime:      c.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:        c.EndTime.Format("2006-01-02 15:04:05"),
		CreatedAt:      c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
