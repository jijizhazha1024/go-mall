package code

const (
	// auth
	AuthBlank = 10000 + iota
	AuthExpired
	AuthSuccess
	AuthFail
	TokenRenewed
	TokenRenewalFailed
	TokenValid
	TokenInvalid
	IllegalProxyAddress
	AuthExpiredByLogout
)
const (
	AuthBlankMsg           = "认证信息为空"
	AuthExpiredMsg         = "认证过期或不存在"
	AuthSuccessMsg         = "身份令牌分发成功"
	AuthFailMsg            = "身份令牌分发失败"
	TokenRenewedMsg        = "令牌续期成功"
	TokenRenewalFailedMsg  = "令牌续期失败"
	TokenValidMsg          = "令牌有效"
	TokenInvalidMsg        = "令牌无效"
	IllegalProxyAddressMsg = "非法代理地址"
	AuthExpiredByLogoutMsg = "认证过期，用户已注销"
)
