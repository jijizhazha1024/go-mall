package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/apis/product/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
	"jijizhazha1024/go-mall/services/product/productcatalogservice"
)

type ServiceContext struct {
	Config                config.Config
	ProductRpc            productcatalogservice.ProductCatalogService
	WithClientMiddleware  rest.Middleware
	WrapperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		ProductRpc:            productcatalogservice.NewProductCatalogService(zrpc.MustNewClient(c.ProductRpc)),
		WithClientMiddleware:  middleware.WithClientMiddleware,
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc, nil, c.OptionPathList),
	}
}
