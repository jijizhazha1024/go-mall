// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"jijizhazha1024/go-mall/apis/ai/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.WithClientMiddleware,serverCtx.WrapperAuthMiddleware},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/",
					Handler: AiHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/douyin/ai"),
	)
}
