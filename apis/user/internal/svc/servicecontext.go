package svc

import (
	"jijizhazha1024/go-mall/apis/user/internal/config"
	"jijizhazha1024/go-mall/common/middleware"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                config.Config
	WrapperAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.UserRpc), // # 需要指定认证rpc地址
	}
}
