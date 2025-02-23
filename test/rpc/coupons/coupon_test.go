package coupons

import (
	"context"
	"encoding/json"
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
	var price = 299900
	productID := 3

	t.Run("折扣优惠券", func(t *testing.T) {
		disCount := 80
		quanity := 2
		discountAmount := price * quanity * (100 - disCount) / 100
		final := price*quanity - discountAmount
		//	ZK20250214001
		coupon, err2 := couponsClient.CalculateCoupon(context.Background(), &coupons.CalculateCouponReq{
			CouponId: "ZK20250214001",
			UserId:   1,
			Items: []*coupons.Items{
				{
					ProductId: int32(productID),
					Quantity:  int32(quanity),
				},
			},
		})
		if err2 != nil {
			t.Error(err2)
			return
		}
		assert.Equal(t, uint32(0), coupon.StatusCode)
		assert.Equal(t, final, coupon.FinalAmount)

	})
}
