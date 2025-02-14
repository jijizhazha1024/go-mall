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
		return
	}
	assert.Equal(t, uint32(0), resp.StatusCode)
	for _, coupon := range resp.Coupons {
		t.Log(coupon)
	}
}

// 测试获取优惠券, 优惠券不存在
func Test_GetCouponLogic_GetCoupon_NotFount(t *testing.T) {

	resp, err := couponsClient.GetCoupon(context.Background(), &coupons.GetCouponReq{
		Id: "1",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != 0 {
		t.Logf("code：%d, msg:%s", resp.StatusCode, resp.StatusMsg)
		return
	}

}
func Test_GetCouponLogic_GetCoupon(t *testing.T) {

	resp, err := couponsClient.GetCoupon(context.Background(), &coupons.GetCouponReq{
		Id: "67508ec1ea7111ef86d80242ac120005",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, uint32(0), resp.StatusCode)
	t.Log(resp.Coupon)
}
