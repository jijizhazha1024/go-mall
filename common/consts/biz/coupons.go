package biz

import "errors"

const (
	CouponsRpcPort = 10009

	UserCouponKey     = "mall:coupon:lock:%d:%s"
	PreOrderCouponKey = "mall:coupon:preorder:%s"
	LockCouponExpire  = 30 * 60 // 30分钟过期（单位：秒）
)

const (
	MaxPageSize     = 30
	DefaultPageSize = 10
)

var (
	CouponsScriptErr = errors.New("coupons script error")
	LockCouponsErr   = errors.New("lock coupons error")
)
