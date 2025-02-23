package code

const (
	OrderNotExist = 11000 + iota
	OrderStatusInvalid
	PaymentStatusInvalid
	OrderExist
	OrderExpired
	UserOrderAddressNotExist
	UserOrderItemNotExist
	OrderParameterInvalid
)

const (
	OrderNotExistMsg            = "订单不存在"
	OrderStatusInvalidMsg       = "订单状态无效"
	PaymentStatusInvalidMsg     = "订单支付状态无效"
	OrderExistMsg               = "订单已存在"
	OrderExpiredMsg             = "订单已过期"
	UserOrderAddressNotExistMsg = "用户订单地址不存在"
	UserOrderItemNotExistMsg    = "用户订单关联商品不存在"
	OrderParameterInvalidMsg    = "订单参数无效"
)
