package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CancelOrder 取消订单 由用户发起
func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderRequest) (*order.EmptyRes, error) {
	res := &order.EmptyRes{}
	var orderRes *order2.Orders
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		oRes, err := l.svcCtx.OrderModel.WithSession(session).GetOrderByOrderIDAndUserIDWithLock(ctx, in.OrderId, int32(in.UserId))
		orderRes = oRes
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				res.StatusCode = code.OrderNotExist
				res.StatusMsg = code.OrderNotExistMsg
				return nil
			}
			return err
		}
		switch order.OrderStatus(oRes.OrderStatus) {
		case order.OrderStatus_ORDER_STATUS_CREATED, order.OrderStatus_ORDER_STATUS_PENDING_PAYMENT:
			//	可以取消
			if err := l.svcCtx.OrderModel.WithSession(session).CancelOrder(ctx, int32(in.UserId), in.OrderId, in.CancelReason); err != nil {
				return err
			}
			return nil
		case order.OrderStatus_ORDER_STATUS_PAID:
			res.StatusCode = code.OrderAlreadyPaid
			res.StatusMsg = code.OrderAlreadyPaidMsg
			return nil
		case order.OrderStatus_ORDER_STATUS_COMPLETED:
			res.StatusCode = code.OrderAlreadyCompleted
			res.StatusMsg = code.OrderAlreadyCompletedMsg
			return nil
		case order.OrderStatus_ORDER_STATUS_CANCELLED:
			res.StatusCode = code.OrderAlreadyCancelled
			res.StatusMsg = code.OrderAlreadyCancelledMsg
			return nil
		case order.OrderStatus_ORDER_STATUS_CLOSED:
			res.StatusCode = code.OrderAlreadyClosed
			res.StatusMsg = code.OrderAlreadyClosedMsg
			return nil
		case order.OrderStatus_ORDER_STATUS_REFUND:
			res.StatusCode = code.OrderAlreadyRefund
			res.StatusMsg = code.OrderAlreadyRefundMsg
			return nil

		}
		return nil
	}); err != nil {
		l.Logger.Errorw("cancel order failed", logx.Field("err", err),
			logx.Field("order_id", in.OrderId), logx.Field("user_id", in.UserId))
		return nil, err
	}
	if res.StatusCode != code.Success {
		return res, nil
	}
	// 退还优惠券，退还库存。
	releaseCheckout, err := l.svcCtx.CheckoutRpc.ReleaseCheckout(l.ctx, &checkout.ReleaseReq{
		PreOrderId: orderRes.PreOrderId,
		UserId:     int32(in.UserId),
	})
	if err != nil {
		l.Logger.Errorw("call rpc ReleaseCheckout failed", logx.Field("err", err))
		return nil, err
	}
	if releaseCheckout.StatusCode != code.Success {
		res.StatusCode = releaseCheckout.StatusCode
		res.StatusMsg = releaseCheckout.StatusMsg
		return res, nil
	}
	couponRes, err := l.svcCtx.CouponRpc.ReleaseCoupon(l.ctx, &coupons.ReleaseCouponReq{
		PreOrderId:   orderRes.PreOrderId,
		UserId:       int32(in.UserId),
		Reason:       in.CancelReason,
		UserCouponId: orderRes.CouponId,
	})
	if err != nil {
		l.Logger.Errorw("call rpc ReleaseCoupon failed", logx.Field("err", err))
		return nil, err
	}
	if couponRes.StatusCode != code.Success {
		res.StatusCode = couponRes.StatusCode
		res.StatusMsg = couponRes.StatusMsg
		return res, nil
	}
	orderItems, err := l.svcCtx.OrderItemModel.QueryOrderItemsByOrderID(l.ctx, orderRes.OrderId)
	if err != nil {
		l.Logger.Errorw("query order items failed", logx.Field("err", err))
		return nil, err
	}
	inventoryItems := make([]*inventory.InventoryReq_Items, len(orderItems))
	for i, item := range orderItems {
		inventoryItems[i] = &inventory.InventoryReq_Items{
			ProductId: int32(item.ProductId),
			Quantity:  int32(item.Quantity),
		}
	}
	// 退还库存
	inventoryResp, err := l.svcCtx.InventoryRpc.ReturnInventory(l.ctx, &inventory.InventoryReq{
		PreOrderId: orderRes.PreOrderId,
		UserId:     int32(in.UserId),
		Items:      inventoryItems,
	})
	if err != nil {
		l.Logger.Errorw("call rpc ReturnInventory failed", logx.Field("err", err))
		return nil, err
	}
	if inventoryResp.StatusCode != code.Success {
		res.StatusCode = inventoryResp.StatusCode
		res.StatusMsg = inventoryResp.StatusMsg
	}
	return res, nil
}
