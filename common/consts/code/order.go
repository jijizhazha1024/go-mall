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
	OrderAlreadyPaid
	OrderAlreadyCompleted
	OrderAlreadyCancelled
	OrderAlreadyClosed
	OrderAlreadyRefund
	CreateOrderFailed
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
	OrderAlreadyPaidMsg         = "订单已支付，不可再修改"
	OrderAlreadyCompletedMsg    = "订单已完成，不可再修改"
	OrderAlreadyCancelledMsg    = "订单已取消，请勿重复操作"
	OrderAlreadyClosedMsg       = "订单已关闭，请勿重复操作"
	OrderAlreadyRefundMsg       = "订单已退款，请勿重复操作"
	CreateOrderFailedMsg        = "创建订单失败"
)
