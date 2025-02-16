package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/coupons/coupon"
	"jijizhazha1024/go-mall/dal/model/coupons/user_coupons"

	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClaimCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCouponLogic {
	return &ClaimCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClaimCoupon 用户领取优惠券
func (l *ClaimCouponLogic) ClaimCoupon(in *coupons.ClaimCouponReq) (*coupons.ClaimCouponResp, error) {

	res := &coupons.ClaimCouponResp{}
	var couponQuery *coupon.Coupons
	if in.CouponId == "" {
		res.StatusCode = code.NotWithParam
		res.StatusMsg = code.NotWithParamMsg
		return res, nil
	}
	// --------------- 为用户创建优惠券 事务 ---------------
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

		// --------------- check ---------------
		//  check user already claimed
		exist, err := l.svcCtx.UserCouponsModel.CheckUserCouponExistWithSession(ctx, session, uint64(in.UserId), in.CouponId)
		if err != nil {
			logx.Errorw("query user coupons error", logx.Field("err", err))
			return err
		}
		if exist {
			res.StatusCode = code.CouponsAlreadyClaimed
			res.StatusMsg = code.CouponsAlreadyClaimedMsg
			return nil
		}

		// check coupon stock and status
		one, err := l.svcCtx.CouponsModel.FindOneWithLock(ctx, session, in.CouponId)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				res.StatusCode = code.CouponsNotExist
				res.StatusMsg = code.CouponsNotExistMsg
				return nil
			}
			logx.Errorw("query coupons by id error", logx.Field("err", err))
			return err
		}
		if one.RemainingCount <= 0 {
			res.StatusCode = code.CouponsOutOfStock
			res.StatusMsg = code.CouponsOutOfStockMsg
			return nil
		}
		if one.Status == 0 {
			res.StatusCode = code.CouponsNotAvailable
			res.StatusMsg = code.CouponsNotAvailableMsg
			return nil
		}
		couponQuery = one
		// --------------- create ---------------
		// decrease coupons stock
		if err := l.svcCtx.CouponsModel.DecreaseStockWithSession(ctx, session, in.CouponId, 1); err != nil {
			logx.Errorw("decrease coupons stock error", logx.Field("err", err))
			return err
		}
		// create user coupons
		if _, err := l.svcCtx.UserCouponsModel.WithSession(session).Insert(ctx, &user_coupons.UserCoupons{
			UserId:   uint64(in.UserId),
			CouponId: in.CouponId,
			Status:   int64(coupons.CouponStatus_COUPON_STATUS_UNSPECIFIED),
		}); err != nil {
			logx.Errorw("create user coupons error", logx.Field("err", err))
			return err
		}
		//return errors.New("test rollback")
		return nil
	}); err != nil {
		res.StatusCode = code.ServerError
		res.StatusMsg = code.ServerErrorMsg
		logx.Errorw("create user coupons error", logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("coupon_id", in.CouponId))
		return res, err
	}
	if res.StatusCode != code.Success {
		return res, nil
	}
	res.Coupon = convertCoupon2Resp(couponQuery)
	return res, nil
}
