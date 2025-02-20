package biz

const (
	CouponsRpcPort = 10009

	UserCouponKey     = "user_coupon:%d:%s"
	PreOrderCouponKey = "preorder_coupons:%s"
	LockCouponExpire  = 30 * 60 // 30分钟过期（单位：秒）
)

const (
	MaxPageSize     = 30
	DefaultPageSize = 10
)
