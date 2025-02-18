package coupons

import (
	"context"
	"github.com/stretchr/testify/assert"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"testing"
)

// --------------- 用户优惠卷列表 ---------------
func Test_ListUserCouponsLogic_ListUserCoupons(t *testing.T) {
	userCoupons, err := couponsClient.ListUserCoupons(context.Background(), &coupons.ListUserCouponsReq{
		Pagination: &coupons.PaginationReq{
			Size: 10,
			Page: 1,
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

// --------------- 用户领取优惠卷 ---------------

// 用户领取
func Test_ClaimCouponLogic_ClaimCoupon(t *testing.T) {
	res, err := couponsClient.ClaimCoupon(context.Background(), &coupons.ClaimCouponReq{
		UserId:   1,
		CouponId: "67756077ea7111ef86d80242ac120005",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

// 用户已领取
func Test_ClaimCouponLogic_ClaimCoupon_AlreadyClaimed(t *testing.T) {
	res, err := couponsClient.ClaimCoupon(context.Background(), &coupons.ClaimCouponReq{
		UserId:   1,
		CouponId: "67508ec1ea7111ef86d80242ac120005",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

// 优惠券已售罄
func Test_ClaimCouponLogic_ClaimCoupon_OutOfStock(t *testing.T) {
	res, err := couponsClient.ClaimCoupon(context.Background(), &coupons.ClaimCouponReq{
		UserId:   1,
		CouponId: "679e623cea7111ef86d80242ac120005",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

// 回滚
func Test_ClaimCouponLogic_ClaimCoupon_Rollback(t *testing.T) {
	res, err := couponsClient.ClaimCoupon(context.Background(), &coupons.ClaimCouponReq{
		UserId:   1,
		CouponId: "679e623cea7111ef86d80242ac120005",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

// --------------- 优惠券使用记录 ---------------
func Test_ListCouponUsagesLogic_ListCouponUsages(t *testing.T) {
	couponUsages, err := couponsClient.ListCouponUsages(context.Background(), &coupons.ListCouponUsagesReq{
		Pagination: &coupons.PaginationReq{
			Size: 10,
			Page: 1,
		},
		UserId: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, uint32(0), couponUsages.StatusCode)
	for _, couponUsage := range couponUsages.Usages {
		t.Log(couponUsage)
	}
}
