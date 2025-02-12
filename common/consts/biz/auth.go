package biz

import (
	"time"
)

type CtxKey string

const (
	AuthsRpcPort        = 10000
	UserIDKey    CtxKey = "user_id"
	ClientIPKey  CtxKey = "client_ip"

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
