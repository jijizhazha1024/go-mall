package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/coupons/coupon_usage"
	"time"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UseCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseCouponLogic {
	return &UseCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UseCoupon 使用优惠券（支付成功确认）
func (l *UseCouponLogic) UseCoupon(in *coupons.UseCouponReq) (*coupons.EmptyResp, error) {
	res := &coupons.EmptyResp{}
	// 修改用户优惠券状态，记录优惠券使用记录
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

		// --------------- check ---------------
		// 判断用户优惠券状态是否已经是已使用，支付成功后，修改优惠券状态为已使用
		status, err := l.svcCtx.UserCouponsModel.WithSession(session).GetStatusByUserIdCouponId(ctx, in.UserId, in.CouponId)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				l.Logger.Infow("user coupon not exist", logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId))
				res.StatusCode = code.CouponsNotExist
				res.StatusMsg = code.CouponsNotExistMsg
				return nil
			}
			l.Logger.Errorw("get user coupon status error", logx.Field("err", err),
				logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId),
				logx.Field("order_id", in.OrderId), logx.Field("pre_order_id", in.PreOrderId))
			return err
		}
		// 2. 状态校验
		if coupons.CouponStatus(status.Status) != coupons.CouponStatus_COUPON_STATUS_LOCKED {
			res.StatusCode = code.CouponStatusInvalid
			res.StatusMsg = code.CouponStatusInvalidMsg
			l.Logger.Infow("coupon status invalid", logx.Field("user_id", in.UserId),
				logx.Field("coupon_id", in.CouponId), logx.Field("order_id", in.OrderId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("status", status.Status))
			return nil
		}

		// --------------- query ---------------
		tp, err := l.svcCtx.CouponsModel.GetCouponTypeByID(ctx, session, in.CouponId)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				l.Logger.Infow("coupon not exist", logx.Field("coupon_id", in.CouponId))
				res.StatusCode = code.CouponsNotExist
				res.StatusMsg = code.CouponsNotExistMsg
				return nil
			}
			l.Logger.Errorw("get coupon error", logx.Field("err", err),
				logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId),
				logx.Field("order_id", in.OrderId), logx.Field("pre_order_id", in.PreOrderId))
			return err
		}

		// --------------- update and record ---------------
		// update
		if err := l.svcCtx.UserCouponsModel.WithSession(session).UpdateStatusOrderById(ctx,
			in.OrderId, int(status.ID), coupons.CouponStatus_COUPON_STATUS_USED); err != nil {
			l.Logger.Errorw("update user coupon status error", logx.Field("err", err),
				logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId),
				logx.Field("order_id", in.OrderId), logx.Field("pre_order_id", in.PreOrderId))
			return err
		}
		// record
		if _, err := l.svcCtx.CouponUsageModel.WithSession(session).Insert(ctx, &coupon_usage.CouponUsage{
			OrderId:        in.OrderId,
			CouponId:       in.CouponId,
			UserId:         uint64(in.UserId),
			CouponType:     tp,
			DiscountAmount: in.DiscountAmount,
			OriginValue:    in.OriginAmount,
			AppliedAt:      time.Now(),
		}); err != nil {
			l.Logger.Errorw("insert coupon usage error", logx.Field("err", err),
				logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId),
				logx.Field("order_id", in.OrderId), logx.Field("pre_order_id", in.PreOrderId))
			return err
		}
		return nil
	}); err != nil {
		l.Logger.Errorw("insert coupon usage error", logx.Field("err", err),
			logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId),
			logx.Field("order_id", in.OrderId), logx.Field("pre_order_id", in.PreOrderId))
		return nil, err
	}

	l.Logger.Infow("use coupon success", logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId),
		logx.Field("order_id", in.OrderId), logx.Field("pre_order_id", in.PreOrderId))
	return res, nil
}
