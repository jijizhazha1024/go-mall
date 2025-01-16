package code

// 状态码
const (
	Success     = 200
	Fail        = 400
	ServerError = 500

	// auth

	AuthBlank   = 10000
	AuthExpired = 10001
)

// 状态码描述
const (
	SuccessMsg     = "success"
	FailMsg        = "请求参数错误"
	ServerErrorMsg = "服务器错误"

	AuthBlankMsg   = "认证信息为空"
	AuthExpiredMsg = "认证过期或不存在"
)
