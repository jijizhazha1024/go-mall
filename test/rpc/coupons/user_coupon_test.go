package coupons

import (
	"context"
	"github.com/stretchr/testify/assert"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"testing"
)

func Test_ListUserCouponsLogic_ListUserCoupons(t *testing.T) {
	userCoupons, err := couponsClient.ListUserCoupons(context.Background(), &coupons.ListUserCouponsReq{
		Pagination: &coupons.PaginationReq{
			Limit: 10,
			Page:  1,
		},
		UserId: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, uint32(0), userCoupons.StatusCode)
	for _, coupon := range userCoupons.UserCoupons {
		t.Log(coupon)
	}
}
