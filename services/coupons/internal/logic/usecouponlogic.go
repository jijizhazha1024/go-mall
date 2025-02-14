package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/coupons/coupon_usage"
	"jijizhazha1024/go-mall/dal/model/coupons/user_coupons"
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

// UseCoupon 使用优惠券
func (l *UseCouponLogic) UseCoupon(in *coupons.UseCouponReq) (*coupons.UseCouponResp, error) {

	res := &coupons.UseCouponResp{}
	now := time.Now()
	err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// --------------- check ---------------

		// check orderId
		// 判断订单是否存在，&& 未支付状态 后续添加

		// check coupon state
		one, err := l.svcCtx.CouponsModel.WithSession(session).FindOne(ctx, in.CouponId)
		if err != nil {
			// not exist
			if errors.Is(err, sqlx.ErrNotFound) {
				res.StatusCode = code.CouponsNotExist
				res.StatusMsg = code.CouponsNotExistMsg
			}
			logx.Errorw("query coupons by id error", logx.Field("err", err))
			return err
		}
		if one.Status == 0 {
			res.StatusCode = code.CouponsNotExist
			res.StatusMsg = code.CouponsNotExistMsg
			return nil
		}
		if now.Before(one.StartTime) {
			res.StatusCode = code.CouponsNotStart
			res.StatusMsg = code.CouponsNotStartMsg
			return nil
		}
		if now.After(one.EndTime) {
			res.StatusCode = code.CouponsExpired
			res.StatusMsg = code.CouponsExpiredMsg
			return nil
		}

		// check user coupon
		userCoupon, err := l.svcCtx.UserCouponsModel.GetUserCouponByUserIdCouponIdLock(ctx, session, uint64(in.UserId), in.CouponId)
		if err != nil {
			// not exist
			if errors.Is(err, sqlx.ErrNotFound) {
				res.StatusCode = code.UserNotHaveCoupons
				res.StatusMsg = code.UserNotHaveCouponsMsg
				return nil
			}
			logx.Errorw("query user coupons error", logx.Field("err", err))
			return err
		}
		// check status
		switch coupons.CouponStatus(userCoupon.Status) {
		case coupons.CouponStatus_COUPON_STATUS_USED:
			res.StatusCode = code.CouponsAlreadyUsed
			res.StatusMsg = code.CouponsAlreadyUsedMsg
			return nil
		case coupons.CouponStatus_COUPON_STATUS_EXPIRED:
			res.StatusCode = code.CouponsExpired
			res.StatusMsg = code.CouponsExpiredMsg
			return nil
		}

		// --------------- update ---------------

		if err := l.svcCtx.UserCouponsModel.WithSession(session).Update(ctx, &user_coupons.UserCoupons{
			Id:        userCoupon.Id,
			Status:    int64(coupons.CouponStatus_COUPON_STATUS_USED),
			UsedAt:    sql.NullTime{Time: now, Valid: true},
			OrderId:   sql.NullString{String: in.OrderId, Valid: true},
			UpdatedAt: now,
		}); err != nil {
			logx.Errorw("update user coupons error", logx.Field("err", err))
			return err
		}
		if _, err := l.svcCtx.CouponUsageModel.WithSession(session).Insert(ctx, &coupon_usage.CouponUsage{
			UserId:     uint64(in.UserId),
			CouponId:   in.CouponId,
			OrderId:    in.OrderId,
			CouponType: one.Type,
			AppliedAt:  now,
		}); err != nil {
			logx.Errorw("insert coupon usage error", logx.Field("err", err))
			return err
		}
		// record usage

		return nil
	})
	if err != nil {
		res.StatusCode = code.ServerError
		res.StatusMsg = code.ServerErrorMsg
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId))
		return res, err
	}
	return &coupons.UseCouponResp{}, nil
}
