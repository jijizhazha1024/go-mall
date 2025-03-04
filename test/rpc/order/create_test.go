package order

import (
	"context"
	"jijizhazha1024/go-mall/services/order/order"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	createOrder, err := orderClient.CreateOrder(context.TODO(), &order.CreateOrderRequest{
		PreOrderId:    "019555d7-8dca-7f17-b945-cee24c0efb7b",
		UserId:        1,
		AddressId:     1,
		PaymentMethod: order.PaymentMethod_ALIPAY,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(createOrder)
}
