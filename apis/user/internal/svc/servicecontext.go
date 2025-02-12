package svc

import (
	"jijizhazha1024/go-mall/apis/user/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
	"jijizhazha1024/go-mall/services/auths/authsclient"
	"jijizhazha1024/go-mall/services/users/usersclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	UserRpc               usersclient.Users
	AuthsRpc              authsclient.Auths
	WrapperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		UserRpc:               usersclient.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		AuthsRpc:              authsclient.NewAuths(zrpc.MustNewClient(c.AuthsRpc)),
		Config:                c,
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc),
	}
}
