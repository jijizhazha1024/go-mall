package biz

import (
	"time"
)

const (
	AuthsRpcPort = 10000
	AuthParamKey = "user_id"

	TokenExpire        = time.Second * 10
	TokenRenewalExpire = time.Hour * 24 * 7
)

var (

	// WhitePath 白名单
	WhitePath = []string{
		"/douyin/user/register",
		"/douyin/user/login",
	}
)
