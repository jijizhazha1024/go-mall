package coupons

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"testing"
)

func Test_LockCouponLogic_LockCoupon(t *testing.T) {
	uci := uuid.New().String()[:8]
	pid := uuid.New().String()[:8]
	t.Run("正常情况", func(t *testing.T) {
		res, err := couponsClient.LockCoupon(context.Background(), &coupons.LockCouponReq{
			UserId:       1,
			UserCouponId: uci,
			PreOrderId:   pid,
		})
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, int32(code.Success), res.StatusCode)
	})
	t.Run("优惠卷已经锁定", func(t *testing.T) {
		res, err := couponsClient.LockCoupon(context.Background(), &coupons.LockCouponReq{
			UserId:       1,
			UserCouponId: uci,
			PreOrderId:   pid,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res)
		lock, err := couponsClient.LockCoupon(context.Background(), &coupons.LockCouponReq{
			UserId:       1,
			UserCouponId: uci,
			PreOrderId:   pid,
		})
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, int32(code.CouponsAlreadyLocked), lock.StatusCode)

	})
}

func Test_UnlockCouponLogic_UnlockCoupon(t *testing.T) {
	uci := uuid.New().String()[:8]
	pid := uuid.New().String()[:8]
	t.Run("正常情况", func(t *testing.T) {
		res, err := couponsClient.LockCoupon(context.Background(), &coupons.LockCouponReq{
			UserId:       1,
			UserCouponId: uci,
			PreOrderId:   pid,
		})
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, int32(code.Success), res.StatusCode)
		unlock, err := couponsClient.ReleaseCoupon(context.Background(), &coupons.ReleaseCouponReq{
			UserId:       1,
			UserCouponId: uci,
			PreOrderId:   pid,
		})
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, int32(code.Success), unlock.StatusCode)
	})
	t.Run("优惠卷已经释放", func(t *testing.T) {
		res, err := couponsClient.ReleaseCoupon(context.Background(), &coupons.ReleaseCouponReq{
			UserId:       1,
			UserCouponId: uci,
			PreOrderId:   pid,
		})
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, int32(code.CouponsAlreadyReleased), res.StatusCode)
	})
}

// 用户优惠券使用情况
func Test_ListCouponsUsageLogic_ListCouponsUsage(t *testing.T) {
	res, err := couponsClient.ListCouponUsages(context.Background(), &coupons.ListCouponUsagesReq{
		Pagination: &coupons.PaginationReq{
			Page: 1,
			Size: 10,
		},
		UserId: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Usages)
}
