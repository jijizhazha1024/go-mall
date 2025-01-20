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

func WrapperAuthMiddleware(rpcConf zrpc.RpcClientConf) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// white path
			if slices.Contains(biz.WhitePath, r.URL.Path) {
				next(w, r)
				return
			}
			var token = r.PostFormValue("token")
			if r.Method == http.MethodGet {
				token = r.FormValue("token")
			} else if token == "" && r.Method == http.MethodPost {
				token = r.URL.Query().Get("token")
			}
			// optional token
			matched, err := regexp.MatchString(`^/douyin/products(/.*)?$`, r.URL.Path)
			if err != nil {
				logx.Errorw("regexp match failed", logx.Field("err", err))
				httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.ServerError, code.ServerErrorMsg))
				return
			}
			if token == "" && matched {
				logx.Infof("token is blank")
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
				logx.Infow("auth failed", logx.Field("status_msg", res.StatusMsg))
				httpx.OkJsonCtx(r.Context(), w, response.NewResponse(int(res.StatusCode), res.StatusMsg))
				return
			}
			// with user_id, 后面都可以在 请求中获取user_id
			r.Form.Set(biz.AuthParamKey, strconv.Itoa(int(res.UserId)))
			next(w, r)
		}
	}
}
