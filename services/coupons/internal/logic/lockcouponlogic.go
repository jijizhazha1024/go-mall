package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/svc"
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
func (l *LockCouponLogic) LockCoupon(in *coupons.LockCouponReq) (*coupons.EmptyResp, error) {
	res := &coupons.EmptyResp{}

	// --------------- check ---------------
	if in.UserId == 0 || len(in.UserCouponId) == 0 || len(in.PreOrderId) == 0 {
		res.StatusCode = code.NotWithParam
		res.StatusMsg = code.NotWithParamMsg
		return nil, status.Error(codes.Aborted, code.NotWithParamMsg)
	}

	// --------------- transact ---------------
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1.检查优惠券状态
		expired, err := l.svcCtx.CouponsModel.CheckExpirationAndStatus(l.ctx, session, in.UserCouponId)
		if err != nil {
			logx.Errorw("check coupon status error", logx.Field("err", err))
			return err
		}
		if expired {
			res.StatusCode = code.CouponsExpired
			res.StatusMsg = code.CouponsExpiredMsg
			return nil
		}
		// 2. 用户是否有该优惠券
		exist, err := l.svcCtx.UserCouponsModel.CheckUserCouponExistWithLock(l.ctx, session, uint64(in.UserId), in.UserCouponId)
		if err != nil {
			logx.Errorw("check user coupon exist error", logx.Field("err", err))
			return err
		}
		if !exist {
			res.StatusCode = code.UserNotHaveCoupons
			res.StatusMsg = code.UserNotHaveCouponsMsg
			return nil
		}
		if err := l.svcCtx.UserCouponsModel.LockUserCoupon(l.ctx, session, uint64(in.UserId), in.UserCouponId); err != nil {
			logx.Errorw("update coupon status error", logx.Field("err", err))
			return err
		}
		return nil
	}); err != nil {
		logx.Errorw("transact lock coupon error", logx.Field("err", err))
		// !!!一般数据库不会错误不需要dtm回滚，就让他一直重试
		return nil, status.Error(codes.Internal, code.ServerErrorMsg) // 触发重试
	}
	if res.StatusCode != 0 {
		return nil, status.Error(codes.Aborted, res.StatusMsg)
	}
	return &coupons.EmptyResp{}, nil
}
