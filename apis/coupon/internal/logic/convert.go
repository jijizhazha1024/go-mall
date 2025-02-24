package logic

import (
	"jijizhazha1024/go-mall/apis/coupon/internal/types"
	"jijizhazha1024/go-mall/services/coupons/coupons"
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

func convertToCouponItems(items []types.Items) []*coupons.Items {
	couponItems := make([]*coupons.Items, len(items))
	for i, item := range items {
		couponItems[i] = &coupons.Items{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	return couponItems
}
func convertCouponUsageList2Resp(usages []*coupons.CouponUsage) []types.CouponUsage {
	res := make([]types.CouponUsage, 0, len(usages))
	for i, usage := range usages {
		res[i] = types.CouponUsage{
			ID:         usage.Id,
			UserID:     usage.UserId,
			OrderID:    usage.OrderId,
			PreOrderID: usage.PreOrderId,
			CouponID:   usage.CouponId,
			CouponType: usage.CouponType.String(),
			// 确保浮点数精度
			OriginValue:    usage.OriginValue,
			DiscountAmount: usage.DiscountAmount,
			AppliedAt:      usage.AppliedAt,
		}
	}
	return res
}
