package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/response"
	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/auths/authsclient"
)

func WrapperAuthMiddleware(rpcConf zrpc.RpcClientConf) func(next http.HandlerFunc) http.HandlerFunc {
	// 预先生成白名单集合和可选路径前缀
	whitePathSet := make(map[string]struct{})
	for _, path := range biz.WhitePath {
		whitePathSet[path] = struct{}{}
	}
	optionalTokenPrefix := "/douyin/products"

	return func(next http.HandlerFunc) http.HandlerFunc {
		var (
			once    sync.Once
			authRpc authsclient.Auths
		)

		return func(w http.ResponseWriter, r *http.Request) {
			// 白名单路径直接放行
			if _, ok := whitePathSet[r.URL.Path]; ok {
				next(w, r)
				return
			}

			// 获取认证令牌
			token := r.Header.Get(biz.TokenKey)
			refreshToken := r.Header.Get(biz.RefreshTokenKey)

			// 处理可选令牌路径
			isOptionalPath := r.URL.Path == optionalTokenPrefix || strings.HasPrefix(r.URL.Path, optionalTokenPrefix+"/")
			if isOptionalPath && token == "" {
				next(w, r)
				return
			}

			// 非可选路径必须携带令牌
			if token == "" {
				sendAuthError(w, r, code.AuthBlank, code.AuthBlankMsg)
				return
			}

			// 延迟初始化认证客户端
			once.Do(func() {
				authRpc = authsclient.NewAuths(zrpc.MustNewClient(rpcConf))
			})

			// 执行认证流程
			authRes, err := authRpc.Authentication(r.Context(), &auths.AuthReq{Token: token})
			if err != nil {
				logx.Errorw("back err", logx.Field("err", err), logx.Field("token", maskToken(token)), logx.Field("path", r.URL.Path))
				sendServerError(w, r)
				return
			}
			// 处理认证结果
			switch authRes.StatusCode {
			case code.Success:
				setUserContext(r, authRes.UserId)
				next(w, r)
			case code.AuthExpired:
				handleTokenExpiration(w, r, authRpc, refreshToken)
			default:
				sendAuthError(w, r, int(authRes.StatusCode), authRes.StatusMsg)
			}
		}
	}
}

// 设置用户上下文
func setUserContext(r *http.Request, userId uint32) {
	ctx := context.WithValue(r.Context(), biz.AuthParamKey, userId)
	*r = *r.WithContext(ctx)
}

// 处理令牌过期
func handleTokenExpiration(w http.ResponseWriter, r *http.Request, client authsclient.Auths, refreshToken string) {
	if refreshToken == "" {
		sendAuthError(w, r, code.AuthExpired, code.AuthExpiredMsg)
		return
	}

	renewRes, err := client.RenewToken(r.Context(), &auths.AuthRenewalReq{RefreshToken: refreshToken})
	if err != nil {
		logx.Errorw("refresh token err",
			logx.Field("err", err),
			logx.Field("refreshToken", maskToken(refreshToken)),
			logx.Field("path", r.URL.Path))
		sendServerError(w, r)
		return
	}
	if renewRes.StatusCode == code.Success {
		// 由客户端处理刷新结果，进行再次请求
		// 返回状态码和刷新结果
		httpx.OkJsonCtx(r.Context(), w, response.NewRefreshResponse(renewRes))
		return
	}

	sendAuthError(w, r, int(renewRes.StatusCode), renewRes.StatusMsg)
}

func sendAuthError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	httpx.OkJsonCtx(r.Context(), w, response.NewResponse(statusCode, message))
}

func sendServerError(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.ServerError, code.ServerErrorMsg))
}

// 令牌脱敏处理
func maskToken(token string) string {
	if len(token) < 8 {
		return "***"
	}
	return token[:3] + "****" + token[len(token)-3:]
}
