package coupons

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"testing"
)

func Test_ListCouponsLogic_ListCoupons(t *testing.T) {

	resp, err := couponsClient.ListCoupons(context.Background(), &coupons.ListCouponsReq{
		Pagination: &coupons.PaginationReq{
			Page: 1,
			Size: 10,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, uint32(0), resp.StatusCode)
	for _, coupon := range resp.Coupons {
		marshal, _ := json.Marshal(coupon)
		t.Log(string(marshal))
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

func Test_CalculateCouponLogic_CalculateCoupon(t *testing.T) {

	calRes, err := couponsClient.CalculateCoupon(context.Background(), &coupons.CalculateCouponReq{
		CouponId: "LJ20250214001",
		Items: []*coupons.Items{
			{
				ProductId: 1,
				Quantity:  1,
			},
		},
		UserId: 1,
	})
	if err != nil {
		t.Error(err)
	}
	if calRes.StatusCode != 0 {
		t.Logf("code：%d, msg:%s", calRes.StatusCode, calRes.StatusMsg)
		return
	}
	assert.Equal(t, uint32(0), calRes.StatusCode)
	fmt.Println(calRes)
}
