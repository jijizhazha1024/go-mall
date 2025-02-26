package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/code"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrder2PaymentSuccessRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentSuccessRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentSuccessRollbackLogic {
	return &UpdateOrder2PaymentSuccessRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2PaymentSuccessRollback 支付失败的补充操作
func (l *UpdateOrder2PaymentSuccessRollbackLogic) UpdateOrder2PaymentSuccessRollback(in *order.UpdateOrder2PaymentSuccessRequest) (*order.EmptyRes, error) {
	// 参数校验保持与支付成功逻辑一致
	if in.OrderId == "" || in.UserId == 0 {
		return nil, status.Error(codes.Aborted, "参数错误")
	}
	res := &order.EmptyRes{}
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 查询订单并加锁
		orderRes, err := l.svcCtx.OrderModel.WithSession(session).GetOrderByOrderIDAndUserIDWithLock(l.ctx, in.OrderId, in.UserId)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				res.StatusCode = code.OrderNotExist
				res.StatusMsg = code.OrderNotExistMsg
				l.Logger.Infow("order not found",
					logx.Field("order_id", in.OrderId),
					logx.Field("user_id", in.UserId))
				return nil
			}
			l.Logger.Errorw("query order status failed", logx.Field("err", err),
				logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
		}

		// 2. 状态检查（检查是否处于可回滚状态）
		if order.OrderStatus(orderRes.OrderStatus) != order.OrderStatus_ORDER_STATUS_PAID ||
			order.PaymentStatus(orderRes.PaymentStatus) != order.PaymentStatus_PAYMENT_STATUS_PAID {
			res.StatusCode = code.OrderStatusInvalid
			res.StatusMsg = code.OrderStatusInvalidMsg
			l.Logger.Infow("invalid status for rollback",
				logx.Field("order_status", orderRes.OrderStatus),
				logx.Field("payment_status", orderRes.PaymentStatus),
				logx.Field("user_id", in.UserId),
				logx.Field("order_id", in.OrderId))
			return nil
		}

		// 3. 更新订单状态为待支付，支付状态为失败
		if err := l.svcCtx.OrderModel.WithSession(session).UpdateOrder2PaymentRollback(
			l.ctx,
			in.OrderId,
			in.UserId,
		); err != nil {
			l.Logger.Errorw("update order status error", logx.Field("err", err),
				logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
			return err
		}

		return nil
	}); err != nil {
		l.Logger.Errorw("update order status error", logx.Field("err", err),
			logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
		return nil, status.Error(codes.Internal, "更新订单状态失败")
	}
	if res.StatusCode != code.Success {
		l.Logger.Infow("UpdateOrder2PaymentSuccessRollback process aborted",
			logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
		return nil, status.Error(codes.Aborted, res.StatusMsg)
	}
	return &order.EmptyRes{}, nil
}
