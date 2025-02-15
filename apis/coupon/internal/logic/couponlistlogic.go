package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/apis/coupon/internal/svc"
	"jijizhazha1024/go-mall/apis/coupon/internal/types"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponListLogic {
	return &CouponListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponListLogic) CouponList(req *types.CouponListReq) (resp *types.CouponListResp, err error) {

	res, err := l.svcCtx.CouponRpc.ListCoupons(l.ctx, &couponsclient.ListCouponsReq{
		Pagination: &couponsclient.PaginationReq{
			Page:  int32(req.Page),
			Limit: int32(req.PageSize),
		},
		Type: int32(req.Type),
	})
	if err != nil {
		if res != nil && res.StatusCode != code.Success {
			// 处理用户级别info 错误
			return nil, errors.New(int(res.StatusCode), res.StatusMsg)
		}
		l.Logger.Errorf("call rpc ListCoupons failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	resp = &types.CouponListResp{
		CouponList: make([]types.CouponItemResp, 0, len(res.Coupons)),
	}
	for _, item := range res.Coupons {
		resp.CouponList = append(resp.CouponList, *convertCoupon2Resp(item))
	}
	return
}
