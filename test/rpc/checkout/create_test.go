package checkout

import (
	"context"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"testing"
)

func TestPrePareCheckout(t *testing.T) {
	resp, err := checkoutClient.PrepareCheckout(context.TODO(), &checkout.CheckoutReq{
		UserId: 1,
		OrderItems: []*checkout.CheckoutReq_OrderItem{
			{
				ProductId: 11,
				Quantity:  1,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
