package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"
)

type CreateOrderRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderRollbackLogic {
	return &CreateOrderRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderRollbackLogic) CreateOrderRollback(in *order.CreateOrderRequest) (*order.EmptyRes, error) {
	// 获取订单ID
	orderID, err := l.svcCtx.OrderModel.GetOrderIDByPreID(l.ctx, in.PreOrderId, int32(in.UserId))
	if err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return nil, status.Error(codes.Aborted, "订单不存在")
		}
		l.Logger.Errorw("get order id failed",
			logx.Field("pre_order_id", in.PreOrderId),
			logx.Field("user_id", in.UserId),
			logx.Field("err", err))
		// 重试
		return nil, status.Error(codes.Internal, "回滚失败")
	}
	// 如果订单不存在则直接返回成功（幂等性设计）
	if orderID == "" {
		return nil, status.Error(codes.Aborted, "订单不存在")
	}

	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 删除订单商品项
		if err := l.svcCtx.OrderItemModel.DeleteOrderItemByOrderID(ctx, session, orderID); err != nil {
			l.Logger.Errorw("delete order items failed", logx.Field("orderID", orderID))
			return err
		}
		// 删除订单地址
		if err := l.svcCtx.OrderAddress.DeleteOrderAddressByOrderID(ctx, session, orderID); err != nil {
			l.Logger.Errorw("delete order address failed", logx.Field("orderID", orderID))
			return err
		}
		// 删除主订单
		if err := l.svcCtx.OrderModel.DeleteOrderByOrderID(ctx, session, orderID); err != nil {
			l.Logger.Errorw("delete order failed", logx.Field("orderID", orderID))
			return err
		}
		return nil
	}); err != nil {
		l.Logger.Errorw("rollback failed", logx.Field("err", err))
		return nil, status.Error(codes.Internal, "回滚失败")
	}

	return &order.EmptyRes{}, nil
}
