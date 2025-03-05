package checkout

import (
	"context"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"testing"
)

func TestUpdateStatus(t *testing.T) {
	order, err := checkoutClient.UpdateStatus2Order(context.TODO(), &checkout.UpdateStatusReq{
		PreOrderId: "019554a5-838c-7414-868d-aba6c1d7c6cd",
		UserId:     1,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(order)
}
