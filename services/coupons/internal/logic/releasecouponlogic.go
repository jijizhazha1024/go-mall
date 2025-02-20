package logic

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/internal/lua"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseCouponLogic {
	return &ReleaseCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReleaseCoupon 释放优惠券（订单取消/超时释放）
func (l *ReleaseCouponLogic) ReleaseCoupon(in *coupons.ReleaseCouponReq) (*coupons.EmptyResp, error) {
	res := &coupons.EmptyResp{}

	// 参数校验
	if in.UserId == 0 || len(in.UserCouponId) == 0 || len(in.PreOrderId) == 0 {
		res.StatusCode = code.NotWithParam
		res.StatusMsg = code.NotWithParamMsg
		return res, nil
	}

	// 构造Redis参数
	keys := []string{
		fmt.Sprintf(biz.UserCouponKey, in.UserId, in.UserCouponId),
		fmt.Sprintf(biz.PreOrderCouponKey, in.PreOrderId),
	}
	args := []interface{}{in.UserCouponId}

	// 执行脚本
	result, err := l.svcCtx.Rdb.EvalCtx(l.ctx, lua.UnlockCouponScript, keys, args)
	if err != nil {
		l.Logger.Errorw("release coupon failed",
			logx.Field("error", err),
			logx.Field("user_id", in.UserId),
			logx.Field("user_coupon_id", in.UserCouponId))
		return nil, biz.CouponsScriptErr
	}
	// 处理结果
	resultInt, ok := result.(int64)
	if !ok {
		l.Logger.Errorw("invalid script result type",
			logx.Field("user_id", in.UserId), logx.Field("user_coupon_id", in.UserCouponId))
		return nil, biz.ReleaseCouponsErr
	}
	if resultInt == 1 {
		l.Logger.Infow("coupon already released", logx.Field("user_id", in.UserId), logx.Field("user_coupon_id", in.UserCouponId))
		res.StatusCode = code.CouponsAlreadyReleased
		res.StatusMsg = code.CouponsAlreadyReleasedMsg
		return res, nil
	}

	return res, nil
}
