package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseCheckoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseCheckoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseCheckoutLogic {
	return &ReleaseCheckoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReleaseCheckout UpdateCheckoutStatus2Success 当订单超时，支付超时，支付退款
func (l *ReleaseCheckoutLogic) ReleaseCheckout(in *checkout.ReleaseReq) (*checkout.EmptyResp, error) {
	err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		cacheKey := fmt.Sprintf("checkout:preorder:%d", in.UserId)
		checkoutRecord, err := l.svcCtx.CheckoutModel.FindOneByUserIdAndPreOrderIdWithSession(l.ctx, session, in.UserId, in.PreOrderId)
		if err != nil {
			l.Logger.Errorw("查询结算记录失败",
				logx.Field("err", err),
				logx.Field("pre_order_id", in.PreOrderId))
			return errors.New(code.QuerySettlementRecordFailedMsg)
		}

		if checkoutRecord.Status == int64(checkout.CheckoutStatus_CANCELLED) || checkoutRecord.Status == int64(checkout.CheckoutStatus_EXPIRED) {
			l.Logger.Infof("订单 %s 已经是失效状态，无需更新", in.PreOrderId)
			return nil
		}

		err = l.svcCtx.CheckoutModel.UpdateStatusWithSession(l.ctx, session, int64(checkout.CheckoutStatus_EXPIRED), in.UserId, in.PreOrderId)
		if err != nil {
			l.Logger.Errorw("更新结算状态失败",
				logx.Field("err", err),
				logx.Field("pre_order_id", in.PreOrderId))
			return errors.New(code.UpdateSettlementStatusFailedMsg)
		}

		if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
			l.Logger.Errorw("删除 Redis 锁失败",
				logx.Field("err", err),
				logx.Field("user_id", in.UserId))
		}

		l.Logger.Infof("成功释放订单 %s 并删除结算锁", in.PreOrderId)
		return nil
	})

	if err != nil {
		l.Logger.Errorw("释放结算记录事务失败",
			logx.Field("err", err))
		return nil, err
	}

	return &checkout.EmptyResp{}, nil
}
