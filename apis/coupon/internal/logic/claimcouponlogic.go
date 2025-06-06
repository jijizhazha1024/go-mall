package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"

	"jijizhazha1024/go-mall/apis/coupon/internal/svc"
	"jijizhazha1024/go-mall/apis/coupon/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCouponLogic {
	return &ClaimCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimCouponLogic) ClaimCoupon(req *types.CouponItemReq) (resp *types.CouponItemResp, err error) {

	userID, ok := l.ctx.Value(biz.UserIDKey).(uint32)
	if !ok {
		return nil, errors.New(code.AuthBlank, code.AuthBlankMsg)
	}

	res, err := l.svcCtx.CouponRpc.ClaimCoupon(l.ctx, &couponsclient.ClaimCouponReq{
		UserId:   int32(userID),
		CouponId: req.CouponID,
	})

	if err != nil {
		if res != nil && res.StatusCode != code.Success {
			// 处理用户级别info 错误
			return nil, errors.New(int(res.StatusCode), res.StatusMsg)
		}
		l.Logger.Errorw("call rpc ClaimCoupon failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}
	if res.StatusCode != code.Success {
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = convertCoupon2Resp(res.Coupon)
	return
}
