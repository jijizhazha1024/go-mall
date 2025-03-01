package code

const (
	PaymentExist = 10001+iota
	PaymentMethodNotSupport
)
const (
	PaymentExistMsg = "该订单已存在支付记录"
	PaymentMethodNotSupportMsg = "不支持的支付方式"
)
