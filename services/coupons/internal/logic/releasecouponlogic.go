package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/code"
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

// ReleaseCoupon 释放优惠券（Saga补偿操作）
func (l *ReleaseCouponLogic) ReleaseCoupon(in *coupons.ReleaseCouponReq) (*coupons.EmptyResp, error) {
	// --------------- 参数校验 ---------------
	if in.UserId <= 0 || len(in.UserCouponId) != 36 || len(in.PreOrderId) != 36 {
		return nil, status.Error(codes.Aborted, "参数格式异常") // 触发补偿
	}
	res := &coupons.EmptyResp{}
	// --------------- 事务操作 ---------------
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 检查优惠券锁定状态与订单匹配
		state, err := l.svcCtx.UserCouponsModel.CheckUserCouponStatus(l.ctx, session, uint64(in.UserId), in.UserCouponId)
		if err != nil {
			l.Logger.Errorw("check lock status failed", logx.Field("error", err))
			return err
		}

		// 2. 状态校验（幂等性保障）
		if coupons.CouponUsageStatus(state) != coupons.CouponUsageStatus_COUPON_USAGE_STATUS_LOCKED {
			l.Logger.Infow("coupon status is not locked", logx.Field("userId", in.UserId), logx.Field("couponId", in.UserCouponId))
			res.StatusCode = code.CouponStatusInvalid
			res.StatusMsg = code.CouponStatusInvalidMsg
			return nil
		}

		// 3. 执行状态更新
		if err := l.svcCtx.UserCouponsModel.UpdateStatusOrderById(
			l.ctx,
			"", // 清空ID
			int(in.UserId),
			coupons.CouponUsageStatus_COUPON_USAGE_STATUS_UNUSED,
		); err != nil {
			l.Logger.Errorw("update coupon status failed", logx.Field("error", err))
			return err
		}
		return nil
	}); err != nil {
		l.Logger.Errorw("transact release coupon error", logx.Field("err", err))
		return nil, status.Error(codes.Internal, code.ServerErrorMsg) // 错误已携带正确status
	}
	if res.StatusCode != code.Success {
		return nil, status.Error(codes.Aborted, res.StatusMsg)
	}
	return res, nil
}
