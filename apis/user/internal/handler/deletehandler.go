package handler

import (
	"net/http"

	"jijizhazha1024/go-mall/apis/user/internal/logic"
	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

func DeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDeleteLogic(r.Context(), svcCtx)
		resp, err := l.Delete(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
