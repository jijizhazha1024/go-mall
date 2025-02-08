package biz

import (
	"time"
)

const (
	AuthsRpcPort = 10000
	AuthParamKey = "user_id"

	TokenExpire        = time.Hour * 2
	TokenRenewalExpire = time.Hour * 24 * 7

	TokenKey        = "access_token"
	RefreshTokenKey = "refresh_token"
)

var (
	// WhitePath 白名单
	WhitePath = []string{
		"/douyin/user/register",
		"/douyin/user/login",
	}
)
