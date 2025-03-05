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

// 记录使用优惠券
func Test_UseCouponLogic_UseCoupon(t *testing.T) {
	cid := "FJ20250214001"
	poid := uuid.New().String()
	oid := uuid.New().String()
	uid := 1

	t.Run("正常情况", func(t *testing.T) {
		res, err := couponsClient.UseCoupon(context.Background(), &coupons.UseCouponReq{
			UserId:         1,
			CouponId:       cid,
			OrderId:        oid,
			DiscountAmount: 100,
			OriginAmount:   100,
			PreOrderId:     poid,
		})
		assert.NoError(t, err)
		t.Log(res)
	})
	t.Run("优惠券不存在", func(t *testing.T) {
		invalidCid := "INVALID_CID"
		res, err := couponsClient.UseCoupon(context.Background(), &coupons.UseCouponReq{
			UserId:         int32(uid),
			CouponId:       invalidCid,
			OrderId:        uuid.NewString(),
			DiscountAmount: 100,
			OriginAmount:   100,
		})
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, int32(code.CouponsNotExist), res.StatusCode)
		t.Log(res)
	})
	t.Run("优惠券状态非锁定", func(t *testing.T) {
		// 先设置优惠券为已使用状态
		res, err := couponsClient.UseCoupon(context.Background(), &coupons.UseCouponReq{
			UserId:         int32(uid),
			CouponId:       cid,
			OrderId:        uuid.NewString(),
			DiscountAmount: 100,
			OriginAmount:   100,
		})
		assert.NoError(t, err) // 事务错误处理方式特殊
		assert.Equal(t, int32(code.CouponStatusInvalid), res.StatusCode)

	})
	t.Run("重复使用优惠券", func(t *testing.T) {
		// TODO 将状态改为未使用
		oid := uuid.NewString()
		poid := uuid.New().String()
		// 第一次使用
		res, err := couponsClient.UseCoupon(context.Background(), &coupons.UseCouponReq{
			PreOrderId:     poid,
			UserId:         int32(uid),
			CouponId:       cid,
			OrderId:        oid,
			DiscountAmount: 100,
			OriginAmount:   100,
		})

		assert.NoError(t, err)
		assert.Equal(t, int32(code.Success), res.GetStatusCode())
		// 第二次使用
		res2, err := couponsClient.UseCoupon(context.Background(), &coupons.UseCouponReq{
			UserId:         int32(uid),
			CouponId:       cid,
			OrderId:        oid,
			DiscountAmount: 100,
			OriginAmount:   100,
		})

		assert.NoError(t, err)
		// 第一次使用后状态就变了
		assert.Equal(t, int32(code.CouponStatusInvalid), res2.GetStatusCode())
	})

}
