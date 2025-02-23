package logic

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/lua"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type LockCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LockCouponLogic {
	return &LockCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LockCoupon 锁定优惠券
// LockCoupon 使用Lua脚本实现原子化优惠券锁定
func (l *LockCouponLogic) LockCoupon(in *coupons.LockCouponReq) (*coupons.EmptyResp, error) {

	res := &coupons.EmptyResp{}
	// 参数校验
	if in.UserId == 0 || len(in.UserCouponId) == 0 || len(in.PreOrderId) == 0 {
		res.StatusCode = code.NotWithParam
		res.StatusMsg = code.NotWithParamMsg
		return res, nil
	}

	// 执行Lua脚本
	keys := []string{
		fmt.Sprintf(biz.UserCouponKey, in.UserId, in.UserCouponId), // keys[1],
		fmt.Sprintf(biz.PreOrderCouponKey, in.PreOrderId),
	}
	args := []interface{}{
		in.UserCouponId,
		biz.LockCouponExpire, // 过期时间
		time.Now().Unix(),
	}
	result, err := l.svcCtx.Rdb.EvalCtx(l.ctx, lua.LockCouponScript, keys, args)

	if err != nil {
		l.Logger.Errorw("lock coupon failed", logx.Field("err", err),
			logx.Field("user_id", in.UserId), logx.Field("user_coupon_id", in.UserCouponId))
		return nil, biz.CouponsScriptErr
	}
	resultInt, ok := result.(int64)
	if !ok {
		l.Logger.Errorw("convert result failed", logx.Field("user_id", in.UserId), logx.Field("user_coupon_id", in.UserCouponId))
		return nil, biz.LockCouponsErr
	}
	if resultInt == 1 {
		l.Logger.Infow("coupon already locked", logx.Field("user_id", in.UserId), logx.Field("user_coupon_id", in.UserCouponId))
		res.StatusCode = code.CouponsAlreadyLocked
		res.StatusMsg = code.CouponsAlreadyLockedMsg
		return res, nil
	}
	return &coupons.EmptyResp{}, nil
}
