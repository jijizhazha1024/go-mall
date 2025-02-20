package coupons

import (
	"context"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"testing"
)

func Test_LockCouponLogic_LockCoupon(t *testing.T) {
	res, err := couponsClient.LockCoupon(context.Background(), &coupons.LockCouponReq{
		UserId:       1,
		UserCouponId: "JJZZ2CHTT",
		PreOrderId:   "jsdjfkdjfkdjfkdjfkdjfkdj1fkdjfkdjfkdjfkdjfkdjfkdjfkdjfkdjfkdjfkdjfkd",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
