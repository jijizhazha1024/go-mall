package coupons

import (
	"context"
	"github.com/stretchr/testify/assert"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"testing"
)

func Test_ListCouponsLogic_ListCoupons(t *testing.T) {

	resp, err := couponsClient.ListCoupons(context.Background(), &coupons.ListCouponsReq{
		Pagination: &coupons.PaginationReq{
			Page:  1,
			Limit: 10,
		},
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, uint32(0), resp.StatusCode)
	for _, coupon := range resp.Coupons {
		t.Log(coupon)
	}
}
