package middleware

import (
	"context"
	"errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/response"
	"net"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetIP returns request real ip.
func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}

func WithClientMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置客户端ip，到ctx
		clientIP, err := GetIP(r)
		if err != nil || clientIP == "" {
			httpx.OkJsonCtx(r.Context(), w, response.NewResponse(code.IllegalProxyAddress, code.IllegalProxyAddressMsg))
			return
		}
		ctx := context.WithValue(r.Context(), biz.ClientIPKey, clientIP)
		*r = *r.WithContext(ctx)
		next(w, r)
	}
}
