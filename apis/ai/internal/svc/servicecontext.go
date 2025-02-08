package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"jijizhazha1024/go-mall/apis/ai/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
)

type ServiceContext struct {
	Config                config.Config
	WrapperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc),
	}
}
