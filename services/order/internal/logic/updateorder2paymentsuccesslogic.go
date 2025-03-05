package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/order/internal/mq/notify"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrder2PaymentSuccessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentSuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentSuccessLogic {
	return &UpdateOrder2PaymentSuccessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2PaymentSuccess 修改订单状态为支付成功
func (l *UpdateOrder2PaymentSuccessLogic) UpdateOrder2PaymentSuccess(in *order.UpdateOrder2PaymentSuccessRequest) (*order.EmptyRes, error) {
	// 修改订单状态为支付成功，取消订单超时任务，
	if in.OrderId == "" || in.UserId == 0 || in.PaymentResult == nil {
		return nil, status.Error(codes.Aborted, "参数错误")
	}
	res := &order.EmptyRes{}
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		orderRes, err := l.svcCtx.OrderModel.WithSession(session).GetOrderByOrderIDAndUserIDWithLock(l.ctx, in.OrderId, in.UserId)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				res.StatusCode = code.OrderNotExist
				res.StatusMsg = code.OrderNotExistMsg
				l.Logger.Infow("order not found", logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
				return nil
			}
			l.Logger.Errorw("update order status error", logx.Field("err", err),
				logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
			return err
		}

		// --------------- check ---------------
		// check order is pending payment
		if order.OrderStatus(orderRes.OrderStatus) != order.OrderStatus_ORDER_STATUS_PENDING_PAYMENT {
			res.StatusCode = code.OrderStatusInvalid
			res.StatusMsg = code.OrderStatusInvalidMsg
			l.Logger.Infow("order status error", logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId), logx.Field("order_status", orderRes.OrderStatus))
			return nil
		}
		// check payment status
		if order.PaymentStatus(orderRes.PaymentStatus) != order.PaymentStatus_PAYMENT_STATUS_PAYING {
			res.StatusCode = code.PaymentStatusInvalid
			res.StatusMsg = code.PaymentStatusInvalidMsg
			l.Logger.Infow("payment status error", logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId), logx.Field("payment_status", orderRes.PaymentStatus))
			return nil
		}
		if err := l.svcCtx.OrderModel.WithSession(session).UpdateOrder2Payment(l.ctx,
			in.OrderId, in.UserId, in.PaymentResult, order.OrderStatus_ORDER_STATUS_PAID,
			order.PaymentStatus_PAYMENT_STATUS_PAID); err != nil {
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

	//真实扣减库存，用户消费优惠券，优惠券使用记录
	if res.StatusCode != code.Success {
		l.Logger.Infow("UpdateOrder2PaymentSuccess process aborted",
			logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
		return nil, status.Error(codes.Aborted, res.StatusMsg)
	}
	// 消息队列
	if err := l.svcCtx.OrderNotifyMQ.Product(&notify.OrderNotifyReq{
		OrderId:       in.OrderId,
		UserID:        in.UserId,
		TransactionId: in.PaymentResult.TransactionId,
		PaidAmount:    in.PaymentResult.PaidAmount,
		PaidAt:        in.PaymentResult.PaidAt,
	}); err != nil {
		l.Logger.Errorw("send order notify error", logx.Field("err", err),
			logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
		return nil, status.Error(codes.Internal, "发送订单通知失败")
	}
	return res, nil
}
