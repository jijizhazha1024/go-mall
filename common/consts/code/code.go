package code

// 状态码
const (
	Success     = 200
	Fail        = 400
	ServerError = 500

	// auth
	AuthBlank          = 10000
	AuthExpired        = 10001
	AuthSuccess        = 10002
	AuthFail           = 10003
	TokenRenewed       = 10004
	TokenRenewalFailed = 10005
	TokenValid         = 10006
	TokenInvalid       = 10007

	// 用户服务

	UserCreated             = 20001
	UserCreationFailed      = 20002
	UserAlreadyExists       = 20003
	EmailAlreadyExists      = 20004
	LoginSuccess            = 20005
	LoginFailed             = 20006
	InvalidCredentials      = 20007
	LogoutSuccess           = 20008
	LogoutFailed            = 20009
	UserDeleted             = 20010
	UserDeletionFailed      = 20011
	UserUpdated             = 20012
	UserUpdateFailed        = 20013
	UserInfoRetrieved       = 20014
	UserInfoRetrievalFailed = 20015

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
	SuccessMsg     = "success"
	FailMsg        = "请求参数错误"
	ServerErrorMsg = "服务器错误"

	AuthBlankMsg          = "认证信息为空"
	AuthExpiredMsg        = "认证过期或不存在"
	AuthSuccessMsg        = "身份令牌分发成功"
	AuthFailMsg           = "身份令牌分发失败"
	TokenRenewedMsg       = "令牌续期成功"
	TokenRenewalFailedMsg = "令牌续期失败"
	TokenValidMsg         = "令牌有效"
	TokenInvalidMsg       = "令牌无效"

	UserCreatedMsg             = "用户创建成功"
	UserCreationFailedMsg      = "用户创建失败"
	UserAlreadyExistsMsg       = "用户已存在"
	EmailAlreadyExistsMsg      = "邮箱已存在"
	LoginSuccessMsg            = "登录成功"
	LoginFailedMsg             = "登录失败"
	InvalidCredentialsMsg      = "无效的凭证"
	LogoutSuccessMsg           = "登出成功"
	LogoutFailedMsg            = "登出失败"
	UserDeletedMsg             = "用户删除成功"
	UserDeletionFailedMsg      = "用户删除失败"
	UserUpdatedMsg             = "用户信息更新成功"
	UserUpdateFailedMsg        = "用户信息更新失败"
	UserInfoRetrievedMsg       = "用户身份信息获取成功"
	UserInfoRetrievalFailedMsg = "用户身份信息获取失败"

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
