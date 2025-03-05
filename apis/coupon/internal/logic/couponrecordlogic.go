package logic

import (
	"context"
	xerrors "github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/apis/coupon/internal/svc"
	"jijizhazha1024/go-mall/apis/coupon/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponRecordLogic {
	return &CouponRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponRecordLogic) CouponRecord(req *types.CouponListReq) (resp *types.CouponUsageListResp, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > biz.MaxPageSize {
		req.PageSize = biz.MaxPageSize
	}
	userID, ok := l.ctx.Value(biz.UserIDKey).(uint32)
	if !ok {
		return nil, xerrors.New(code.AuthBlank, code.AuthBlankMsg)
	}
	couponUsages, err := l.svcCtx.CouponRpc.ListCouponUsages(l.ctx, &couponsclient.ListCouponUsagesReq{
		UserId: userID,
		Pagination: &couponsclient.PaginationReq{
			Page: req.Page,
			Size: req.PageSize,
		},
	})
	if err != nil {
		l.Logger.Errorw("call rpc ListCouponUsages failed", logx.Field("err", err))
		return nil, xerrors.New(code.ServerError, code.ServerErrorMsg)
	}
	if couponUsages.StatusCode != code.Success {
		return nil, xerrors.New(int(couponUsages.StatusCode), couponUsages.StatusMsg)
	}
	resp = &types.CouponUsageListResp{
		CouponUsageList: convertCouponUsageList2Resp(couponUsages.Usages),
	}
	return
}
