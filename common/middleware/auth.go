package middleware

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/response"
	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/auths/authsclient"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"sync"
)

var once sync.Once
var authRpc authsclient.Auths
var tokenPattern = regexp.MustCompile(`^/douyin/products(/.*)?$`)

func WrapperAuthMiddleware(rpcConf zrpc.RpcClientConf) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// white path
			if slices.Contains(biz.WhitePath, r.URL.Path) {
				next(w, r)
				return
			}
			// get token from header
			token := r.Header.Get("access_token")
			refreshToken := r.Header.Get("refresh_token")

			// optional token for specific paths
			if tokenPattern.MatchString(r.URL.Path) {
				if token == "" {
					logx.Infof("token is blank for optional path: %s", r.URL.Path)
					next(w, r)
					return
				}
			} else if token == "" {
				httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.AuthBlank, code.AuthBlankMsg))
				return
			}
			// init rpc
			once.Do(func() {
				authRpc = authsclient.NewAuths(zrpc.MustNewClient(rpcConf))
			})

			// auth
			res, err := authRpc.Authentication(r.Context(), &auths.AuthReq{
				Token: token,
			})
			// back err
			if err != nil {
				logx.Errorw("back err", logx.Field("err", err), logx.Field("token", token))
				httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.ServerError, code.ServerErrorMsg))
				return
			}
			// auth failed
			if res.StatusCode != 0 {
				// refresh token
				if res.StatusCode == code.AuthExpired && refreshToken != "" {
					renewRes, err := authRpc.RenewToken(r.Context(), &auths.AuthRenewalReq{
						RefreshToken: refreshToken,
					})
					if err != nil {
						logx.Errorw("refresh token err", logx.Field("err", err))
						httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.ServerError, code.ServerErrorMsg))
						return
					}
					if renewRes.StatusCode == code.AuthExpired {
						// refresh token expired
						httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.AuthExpired, code.AuthExpiredMsg))
						return
					}
					httpx.OkJsonCtx(r.Context(), w, renewRes)
					return
				}
				httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.TokenInvalid, code.TokenInvalidMsg))
				return
			}
			// with user_id, 后面都可以在 请求中获取user_id
			r.Form.Set(biz.AuthParamKey, strconv.Itoa(int(res.UserId)))
			next(w, r)
		}
	}
}
