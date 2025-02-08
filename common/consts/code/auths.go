package code

const (
	// auth
	AuthBlank          = 10000
	AuthExpired        = 10001
	AuthSuccess        = 10002
	AuthFail           = 10003
	TokenRenewed       = 10004
	TokenRenewalFailed = 10005
	TokenValid         = 10006
	TokenInvalid       = 10007
)
const (
	AuthBlankMsg          = "认证信息为空"
	AuthExpiredMsg        = "认证过期或不存在"
	AuthSuccessMsg        = "身份令牌分发成功"
	AuthFailMsg           = "身份令牌分发失败"
	TokenRenewedMsg       = "令牌续期成功"
	TokenRenewalFailedMsg = "令牌续期失败"
	TokenValidMsg         = "令牌有效"
	TokenInvalidMsg       = "令牌无效"
)
