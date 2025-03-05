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

type UpdateOrder2PaymentStatusRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentStatusRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentStatusRollbackLogic {
	return &UpdateOrder2PaymentStatusRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2PaymentStatusRollback 补偿操作：回滚订单状态到CREATED
func (l *UpdateOrder2PaymentStatusRollbackLogic) UpdateOrder2PaymentStatusRollback(in *order.UpdateOrder2PaymentRequest) (*order.EmptyRes, error) {
	res := &order.EmptyRes{}

	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// --------------- 校验当前状态 ---------------
		orderRes, err := l.svcCtx.OrderModel.WithSession(session).
			GetOrderByOrderIDAndUserIDWithLock(ctx, in.OrderId, in.UserId)
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
			return err
		}

		// 只允许回滚PENDING_PAYMENT状态的订单
		if order.OrderStatus(orderRes.OrderStatus) != order.OrderStatus_ORDER_STATUS_PENDING_PAYMENT ||
			order.PaymentStatus(orderRes.PaymentStatus) != order.PaymentStatus_PAYMENT_STATUS_PAYING {
			res.StatusCode = code.OrderStatusInvalid
			res.StatusMsg = code.OrderStatusInvalidMsg
			l.Logger.Infow("invalid status for rollback",
				logx.Field("order_status", orderRes.OrderStatus),
				logx.Field("payment_status", orderRes.PaymentStatus),
				logx.Field("order_id", in.OrderId))
			return nil
		}

		// --------------- 执行状态回滚 ---------------
		if err := l.svcCtx.OrderModel.WithSession(session).UpdateOrderStatusByOrderIDAndUserID(ctx, in.OrderId, in.UserId,
			order.OrderStatus_ORDER_STATUS_CREATED, order.PaymentStatus_PAYMENT_STATUS_NOT_PAID); err != nil {
			l.Logger.Errorw("rollback order status failed",
				logx.Field("error", err),
				logx.Field("order_id", in.OrderId))
			return err
		}

		// --------------- 记录补偿日志 ---------------
		return nil
	}); err != nil {
		l.Logger.Errorw("transaction failed",
			logx.Field("error", err),
			logx.Field("order_id", in.OrderId))
		return nil, status.Error(codes.Internal, "补偿操作失败")
	}
	if res.StatusCode != code.Success {
		return nil, status.Error(codes.Aborted, res.StatusMsg)
	}
	return &order.EmptyRes{}, nil
}
