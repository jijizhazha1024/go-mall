package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/apis/carts/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
	"jijizhazha1024/go-mall/services/carts/cartsclient"
)

type ServiceContext struct {
	Config                config.Config
	CartRpc               cartsclient.Cart
	WithClientMiddleware  rest.Middleware
	WrapperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		CartRpc:               cartsclient.NewCart(zrpc.MustNewClient(c.CartRpc)),
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc, c.WhitePathList, c.OptionPathList),
		WithClientMiddleware:  middleware.WithClientMiddleware,
	}
}
