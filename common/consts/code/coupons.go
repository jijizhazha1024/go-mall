package code

const (
	CouponsNotExist = 90000 + iota
	NotWithParam
	CouponsAlreadyClaimed
	CouponsOutOfStock
	CouponsNotAvailable
	UserNotHaveCoupons
	CouponsExpired
	CouponsAlreadyUsed
	CouponsNotStart
	CouponsAlreadyLocked
)

const (
	CouponsNotExistMsg       = "优惠券不存在"
	NotWithParamMsg          = "缺失必要参数"
	CouponsAlreadyClaimedMsg = "优惠券已被领取"
	CouponsOutOfStockMsg     = "优惠券已售罄"
	CouponsNotAvailableMsg   = "优惠券不可用"
	UserNotHaveCouponsMsg    = "用户未拥有优惠券"
	CouponsExpiredMsg        = "优惠券已过期"
	CouponsAlreadyUsedMsg    = "优惠券已被使用"
	CouponsNotStartMsg       = "优惠券未到使用时间"
	CouponsAlreadyLockedMsg  = "优惠券已被锁定"
)
