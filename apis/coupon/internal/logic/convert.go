package logic

import (
	"jijizhazha1024/go-mall/apis/coupon/internal/types"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"
)

func convertCoupon2Resp(c *couponsclient.Coupon) *types.CouponItemResp {
	return &types.CouponItemResp{
		ID:             c.Id,
		Name:           c.Name,
		Type:           uint8(c.Type),
		Value:          c.Value,
		MinAmount:      c.MinAmount,
		TotalCount:     c.TotalCount,
		RemainingCount: c.RemainingCount,
		StartTime:      c.StartTime,
		EndTime:        c.EndTime,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}

}
