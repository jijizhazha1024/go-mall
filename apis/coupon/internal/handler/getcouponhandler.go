package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jijizhazha1024/go-mall/apis/coupon/internal/logic"
	"jijizhazha1024/go-mall/apis/coupon/internal/svc"
)

func GetCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetCouponLogic(r.Context(), svcCtx)
		resp, err := l.GetCoupon()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
