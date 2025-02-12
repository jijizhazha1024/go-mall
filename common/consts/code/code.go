package code

// 状态码
const (

	// 商品服务

	ProductCreated             = 30001
	ProductCreationFailed      = 30002
	ProductUpdated             = 30003
	ProductUpdateFailed        = 30004
	ProductDeleted             = 30005
	ProductDeletionFailed      = 30006
	ProductInfoRetrieved       = 30007
	ProductInfoRetrievalFailed = 30008

	// 购物车服务

	CartCreated             = 40001
	CartCreationFailed      = 40002
	CartCleared             = 40003
	CartClearFailed         = 40004
	CartInfoRetrieved       = 40005
	CartInfoRetrievalFailed = 40006

	// 订单服务

	OrderCreated            = 50001
	OrderCreationFailed     = 50002
	OrderUpdated            = 50003
	OrderUpdateFailed       = 50004
	OrderCancelled          = 50005
	OrderCancellationFailed = 50006

	// 结算

	SettlementSuccess = 60001
	SettlementFailed  = 60002

	// 支付

	PaymentCancelled          = 70001
	PaymentCancellationFailed = 70002
	PaymentSuccess            = 70003
	PaymentFailed             = 70004

	// AI大模型

	OrderQuerySuccess = 80001
	OrderQueryFailed  = 80002
	AutoOrderSuccess  = 80003
	AutoOrderFailed   = 80004
)

// 状态码描述
const (
	ProductCreatedMsg             = "商品创建成功"
	ProductCreationFailedMsg      = "商品创建失败"
	ProductUpdatedMsg             = "商品信息更新成功"
	ProductUpdateFailedMsg        = "商品信息更新失败"
	ProductDeletedMsg             = "商品删除成功"
	ProductDeletionFailedMsg      = "商品删除失败"
	ProductInfoRetrievedMsg       = "商品信息查询成功"
	ProductInfoRetrievalFailedMsg = "商品信息查询失败"

	CartCreatedMsg             = "购物车创建成功"
	CartCreationFailedMsg      = "购物车创建失败"
	CartClearedMsg             = "购物车清空成功"
	CartClearFailedMsg         = "购物车清空失败"
	CartInfoRetrievedMsg       = "购物车信息获取成功"
	CartInfoRetrievalFailedMsg = "购物车信息获取失败"

	OrderCreatedMsg            = "订单创建成功"
	OrderCreationFailedMsg     = "订单创建失败"
	OrderUpdatedMsg            = "订单信息更新成功"
	OrderUpdateFailedMsg       = "订单信息更新失败"
	OrderCancelledMsg          = "订单取消成功"
	OrderCancellationFailedMsg = "订单取消失败"

	SettlementSuccessMsg = "结算成功"
	SettlementFailedMsg  = "结算失败"

	PaymentCancelledMsg          = "支付取消成功"
	PaymentCancellationFailedMsg = "支付取消失败"
	PaymentSuccessMsg            = "支付成功"
	PaymentFailedMsg             = "支付失败"

	OrderQuerySuccessMsg = "订单查询成功"
	OrderQueryFailedMsg  = "订单查询失败"
	AutoOrderSuccessMsg  = "自动下单成功"
	AutoOrderFailedMsg   = "自动下单失败"
)
