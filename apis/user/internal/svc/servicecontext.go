package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"jijizhazha1024/go-mall/apis/user/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
)

type ServiceContext struct {
	Config         config.Config
	WithMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		WithMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc), // 这里需要进行自定义，执行middleware里的认证中间件
	}
}
