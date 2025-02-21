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

type UpdateOrder2PaymentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentStatusLogic {
	return &UpdateOrder2PaymentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2PaymentStatus 更新订单（支付服务回调使用） 更新为支付中
func (l *UpdateOrder2PaymentStatusLogic) UpdateOrder2PaymentStatus(in *order.UpdateOrder2PaymentRequest) (*order.EmptyRes, error) {
	res := &order.EmptyRes{}
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// --------------- 校验订单状态  ---------------
		s, err := l.svcCtx.OrderModel.WithSession(session).GetOrderStatusByOrderIDAndUserIDWithLock(ctx, in.OrderId, in.UserId)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				res.StatusCode = code.OrderNotExist
				res.StatusMsg = code.OrderNotExistMsg
				l.Logger.Infow("order not found", logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
				return nil
			}
			return err
		}
		if order.OrderStatus(s) != order.OrderStatus_ORDER_STATUS_CREATED {
			res.StatusCode = code.OrderStatusInvalid
			res.StatusMsg = code.OrderStatusInvalidMsg
			l.Logger.Infow("order status error", logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId), logx.Field("order_status", s))
			return nil
		}
		// 修改为待支付
		if err := l.svcCtx.OrderModel.WithSession(session).UpdateOrderStatusByOrderIDAndUserID(ctx,
			in.OrderId, in.UserId, order.OrderStatus_ORDER_STATUS_PENDING_PAYMENT); err != nil {
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

		return nil, status.Error(codes.Aborted, res.StatusMsg)
	}
	return &order.EmptyRes{}, nil
}
