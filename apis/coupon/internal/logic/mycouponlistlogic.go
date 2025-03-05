package logic

import (
	"context"
	xerrors "github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"

	"jijizhazha1024/go-mall/apis/coupon/internal/svc"
	"jijizhazha1024/go-mall/apis/coupon/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyCouponListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMyCouponListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyCouponListLogic {
	return &MyCouponListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyCouponListLogic) MyCouponList(req *types.CouponListReq) (resp *types.CouponMyListResp, err error) {

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

	res, err := l.svcCtx.CouponRpc.ListUserCoupons(l.ctx, &couponsclient.ListUserCouponsReq{
		Pagination: &couponsclient.PaginationReq{
			Page: req.Page,
			Size: req.PageSize,
		}, UserId: int32(userID),
	})
	if err != nil {
		l.Logger.Errorw("call rpc ListUserCoupons failed", logx.Field("err", err))
		return nil, xerrors.New(code.ServerError, code.ServerErrorMsg)
	}
	if res.StatusCode != code.Success {
		return nil, xerrors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.CouponMyListResp{
		CouponList: convertCouponMy2Resp(res.UserCoupons),
	}
	return
}

func convertCouponMy2Resp(coupons []*coupons.UserCoupon) []types.CouponMy {
	resp := make([]types.CouponMy, len(coupons))
	for i, coupon := range coupons {
		if coupon == nil {
			continue
		}
		resp[i] = types.CouponMy{
			CouponID:  coupon.CouponId,
			CreatedAt: coupon.CreatedAt,
			ID:        coupon.Id,
			OrderID:   coupon.OrderId,
			Status:    coupon.Status.String(),
			UpdatedAt: coupon.UpdatedAt,
			UsedAt:    coupon.UsedAt,
		}
	}
	return resp
}
