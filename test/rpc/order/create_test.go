package order

import (
	"context"
	"jijizhazha1024/go-mall/services/order/order"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	createOrder, err := orderClient.CreateOrder(context.TODO(), &order.CreateOrderRequest{
		PreOrderId:    "019554a8-9bd9-706f-bec5-13cbaed6e04e",
		UserId:        1,
		AddressId:     1,
		PaymentMethod: order.PaymentMethod_ALIPAY,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(createOrder)
}
