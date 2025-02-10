package biz

import (
	"time"
)

type UserCtxKey string

const (
	AuthsRpcPort            = 10000
	UserIDKey    UserCtxKey = "user_id"
	ClientIPKey             = "client_ip"

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
