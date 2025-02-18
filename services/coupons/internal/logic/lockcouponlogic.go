package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/biz"
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
	// 参数校验
	if in.UserId == 0 || len(in.UserCouponId) == 0 || len(in.PreOrderId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "参数不合法")
	}

	// 预加载Lua脚本（推荐在服务初始化时加载）

	// 执行Lua脚本
	keys := []string{
		fmt.Sprintf(biz.UserCouponKey, in.UserId, in.UserCouponId),
		fmt.Sprintf(biz.PreOrderCouponKey, in.PreOrderId),
	}
	args := []interface{}{
		in.UserCouponId,
		biz.LockCouponExpire, //
		time.Now().Unix(),
	}
	result, err := l.svcCtx.Rdb.EvalCtx(l.ctx, lua.LockCouponScript, keys, args)
	if err != nil {
		return nil, status.Error(codes.Internal, "锁定优惠券失败")
	}
	resultInt, ok := result.(int64)
	if !ok {
		return nil, status.Error(codes.Internal, "锁定优惠券失败")
	}
	if resultInt == 1 {
		return nil, status.Error(codes.FailedPrecondition, "优惠券已被锁定")
	}
	logx.Infow("user lock coupon success",
		logx.Field("user_id", in.UserId), logx.Field("user_coupon_id", in.UserCouponId), logx.Field("pre_order_id", in.PreOrderId))
	return &coupons.EmptyResp{}, nil
}
