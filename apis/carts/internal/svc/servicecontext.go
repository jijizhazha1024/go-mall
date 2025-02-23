package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/apis/carts/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
	"jijizhazha1024/go-mall/services/carts/cartsclient"
	"jijizhazha1024/go-mall/services/product/productcatalogservice"
)

type ServiceContext struct {
	Config                config.Config
	CartsRpc              cartsclient.Cart
	ProductRpc            productcatalogservice.ProductCatalogService
	WithClientMiddleware  rest.Middleware
	WrapperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		CartsRpc:              cartsclient.NewCart(zrpc.MustNewClient(c.CartsRpc)),
		ProductRpc:            productcatalogservice.NewProductCatalogService(zrpc.MustNewClient(c.ProductRpc)),
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc, nil, nil),
		WithClientMiddleware:  middleware.WithClientMiddleware,
	}
}
