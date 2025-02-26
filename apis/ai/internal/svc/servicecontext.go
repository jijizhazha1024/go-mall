package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/apis/ai/internal/config"
	"jijizhazha1024/go-mall/common/middleware"
	"jijizhazha1024/go-mall/services/ai/aiclient"
)

type ServiceContext struct {
	Config                                      config.Config
	WrapperAuthMiddleware, WithClientMiddleware rest.Middleware
	AiRpc                                       aiclient.Ai
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc, nil, nil),
		WithClientMiddleware:  middleware.WithClientMiddleware,
		AiRpc:                 aiclient.NewAi(zrpc.MustNewClient(c.AIRpc)),
	}
}
