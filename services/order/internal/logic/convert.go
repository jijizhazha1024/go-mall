package logic

import (
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/order/order"
	"time"
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

// --------------- resp ---------------
func convertToOrderResp(orderModelRes *order2.Orders) *order.Order {
	resp := &order.Order{
		OrderId:        orderModelRes.OrderId,
		OrderStatus:    order.OrderStatus(orderModelRes.OrderStatus),
		PaymentStatus:  order.PaymentStatus(orderModelRes.PaymentStatus),
		PaymentMethod:  order.PaymentMethod(orderModelRes.PaymentMethod.Int64),
		OriginalAmount: orderModelRes.OriginalAmount,
		PayableAmount:  orderModelRes.PayableAmount,
		PaidAmount:     orderModelRes.PaidAmount.Int64,
		PaidAt:         orderModelRes.PaidAt.Int64,
		DiscountAmount: orderModelRes.DiscountAmount,
		ExpireTime:     time.Unix(orderModelRes.ExpireTime, 0).Format(time.DateTime),
		CreatedAt:      orderModelRes.CreatedAt.Format(time.DateTime),
		UpdatedAt:      orderModelRes.UpdatedAt.Format(time.DateTime),
		PreOrderId:     orderModelRes.PreOrderId,
		Reason:         orderModelRes.Reason.String,
		TransactionId:  orderModelRes.TransactionId.String,
	}

	return resp
}

func convertToOrderItemResp(orderItems []*order2.OrderItems) []*order.OrderItem {
	resp := make([]*order.OrderItem, len(orderItems))
	for i, item := range orderItems {

		resp[i] = &order.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			UnitPrice:   item.Price,
			Quantity:    item.Quantity,
			ProductDesc: item.ProductDesc,
		}
	}
	return resp
}
func convertToOrderAddressResp(address *order2.OrderAddresses) *order.OrderAddress {
	return &order.OrderAddress{
		AddressId:       address.AddressId,
		RecipientName:   address.RecipientName,
		PhoneNumber:     address.PhoneNumber.String,
		Province:        address.Province.String,
		City:            address.City,
		DetailedAddress: address.DetailedAddress,
		OrderId:         address.OrderId,
		CreatedAt:       address.CreatedAt.Format(time.DateTime),
		UpdatedAt:       address.UpdatedAt.Format(time.DateTime),
	}
}
