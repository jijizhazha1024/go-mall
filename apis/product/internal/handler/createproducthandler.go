package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jijizhazha1024/go-mall/apis/product/internal/logic"
	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"
)

func CreateProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateProductReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateProductLogic(r.Context(), svcCtx)
		resp, err := l.CreateProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
