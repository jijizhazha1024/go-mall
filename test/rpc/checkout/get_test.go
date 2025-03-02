package checkout

import (
	"context"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"testing"
)

func TestGetCheckoutDetail(t *testing.T) {
	detail, err := checkoutClient.GetCheckoutDetail(context.TODO(), &checkout.CheckoutDetailReq{
		PreOrderId: "019554a5-838c-7414-868d-aba6c1d7c6cd",
		UserId:     1,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(detail)
}
func TestGetCheckoutList(t *testing.T) {
	list, err := checkoutClient.GetCheckoutList(context.TODO(), &checkout.CheckoutListReq{
		PageSize: 5,
		Page:     1,
		UserId:   1,
	})
	if err != nil {
		t.Error(err)
	}
	for _, v := range list.Data {
		t.Log(v)
	}
}
