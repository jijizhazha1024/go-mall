package code

const (
	OrderNotExist = 11000 + iota
	OrderStatusInvalid
	PaymentStatusInvalid
)

const (
	OrderNotExistMsg        = "订单不存在"
	OrderStatusInvalidMsg   = "订单状态无效"
	PaymentStatusInvalidMsg = "订单支付状态无效"
)
