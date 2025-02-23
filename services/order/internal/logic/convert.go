package logic

import (
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/coupons/coupons"
)

func convertToCouponItems(items []*checkout.CheckoutItem) []*coupons.Items {
	couponItems := make([]*coupons.Items, 0, len(items))
	for _, item := range items {
		couponItems = append(couponItems, &coupons.Items{
			ProductId: int32(item.ProductId),
			Quantity:  item.Quantity,
		})
	}
	return couponItems
}
func convertToOrderItems(orderID string, items []*checkout.CheckoutItem) []*order2.OrderItems {
	orderItems := make([]*order2.OrderItems, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, &order2.OrderItems{
			OrderId:   orderID,
			ProductId: uint64(item.ProductId),
			Quantity:  uint64(item.Quantity),
		})
	}
	return orderItems
}
