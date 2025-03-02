package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/apis/order/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
	"jijizhazha1024/go-mall/services/order/orderservice"
)

type ServiceContext struct {
	Config                config.Config
	WithClientMiddleware  rest.Middleware
	WrapperAuthMiddleware rest.Middleware
	OrderRpc              orderservice.OrderService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		WithClientMiddleware:  middleware.WithClientMiddleware,
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc, nil, nil),
		OrderRpc:              orderservice.NewOrderService(zrpc.MustNewClient(c.OrderRpc)),
	}
}
