package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
	"jijizhazha1024/go-mall/apis/ai/internal/logic"
	"jijizhazha1024/go-mall/apis/ai/internal/svc"
	"jijizhazha1024/go-mall/apis/ai/internal/types"
)

func AiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}
		l := logic.NewAiLogic(r.Context(), svcCtx)
		resp, err := l.Ai(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}
		xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
	}
}
