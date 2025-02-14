package code

const (
	CouponsNotExist = 90000 + iota
	NotWithParam
	CouponsAlreadyClaimed
	CouponsOutOfStock
	CouponsNotAvailable
)

const (
	CouponsNotExistMsg       = "优惠券不存在"
	NotWithParamMsg          = "缺失必要参数"
	CouponsAlreadyClaimedMsg = "优惠券已被领取"
	CouponsOutOfStockMsg     = "优惠券已售罄"
	CouponsNotAvailableMsg   = "优惠券不可用"
)
