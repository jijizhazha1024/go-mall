package code

const (
	OrderNotExist = 11000 + iota
	OrderStatusInvalid
)

const (
	OrderNotExistMsg      = "订单不存在"
	OrderStatusInvalidMsg = "订单状态无效"
)
