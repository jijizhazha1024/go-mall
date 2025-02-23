// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package handler

import (
	"net/http"

	"jijizhazha1024/go-mall/apis/user/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.WithClientMiddleware, serverCtx.WrapperAuthMiddleware},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/address",
					Handler: AddAddressHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/address",
					Handler: DeleteAddressHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/address",
					Handler: UpdateAddressHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/address",
					Handler: GetAddressHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/address/list",
					Handler: AllAddressListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: DeleteHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/info",
					Handler: GetInfoHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/login",
					Handler: LoginHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/logout",
					Handler: LogoutHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/register",
					Handler: RegisterHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/update",
					Handler: UpdateHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/douyin/user"),
	)
}
