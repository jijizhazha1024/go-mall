package handler

import (
	"github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
	"jijizhazha1024/go-mall/common/consts/code"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jijizhazha1024/go-mall/apis/product/internal/logic"
	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"
)

func GetProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetProductReq
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, errors.New(code.Fail, err.Error()))
			return
		}

		l := logic.NewGetProductLogic(r.Context(), svcCtx)
		resp, err := l.GetProduct(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
