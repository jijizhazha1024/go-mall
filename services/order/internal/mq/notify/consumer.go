package notify

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/order/order"
)

// 处理支付成功的回调
func (a *OrderNotifyMQ) consumer(ctx context.Context) {

	ch, err := a.conn.Channel()
	if err != nil {
		logx.Errorw("Failed to open a channel", logx.Field("err", err))
		return
	}

	results, err := ch.Consume(
		QueueName, // 队列名称
		"",        // 消费者标签
		true,      // 自动确认（ack）
		false,     // 排他性
		false,     // 本地消息
		false,     // 等待确认
		nil,       // 参数
	)
	if err != nil {
		logx.Errorw("Failed to register a consumer", logx.Field("err", err))
	}
	logx.Infow("Starting RabbitMQ consumer...")

	for res := range results {
		logx.Infow("start to consume order", logx.Field("body", string(res.Body)))
		var msg *OrderNotifyReq
		if err := json.Unmarshal(res.Body, &msg); err != nil {
			logx.Errorw("failed to unmarshal message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			if err := res.Reject(false); err != nil {
				logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
			continue
		}
		var orderModelRes *order2.Orders
		// --------------- reverse --------------- 幂等
		if err := a.Model.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
			orderModel := a.OrderModel.WithSession(session)
			orderRes, err := orderModel.GetOrderByOrderIDAndUserIDWithLock(ctx, msg.OrderId, msg.UserID)
			if err != nil {
				return err
			}
			orderModelRes = orderRes
			//只进行处理待付款的订单
			if order.OrderStatus(orderRes.OrderStatus) == order.OrderStatus_ORDER_STATUS_PAID {
				return nil
			}

			if err := orderModel.UpdateOrder2Payment(ctx, msg.OrderId, msg.UserID, &order.PaymentResult{
				PaidAmount:    msg.PaidAmount,
				PaidAt:        msg.PaidAt,
				TransactionId: msg.TransactionId,
			}, order.OrderStatus_ORDER_STATUS_PAID, order.PaymentStatus_PAYMENT_STATUS_PAID); err != nil {
				return err
			}

			return nil
		}); err != nil {
			logx.Errorw("update order status error", logx.Field("err", err),
				logx.Field("order_id", msg.OrderId), logx.Field("user_id", msg.UserID))
			if err := res.Reject(true); err != nil {
				logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
		}

		orderItems, err := a.OrderItemModel.QueryOrderItemsByOrderID(ctx, orderModelRes.OrderId)
		if err != nil {
			logx.Errorw("failed to query order items", logx.Field("err", err), logx.Field("body", string(res.Body)))
			if err := res.Reject(true); err != nil {
				logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
			continue
		}
		ItemsReq := make([]*inventory.InventoryReq_Items, len(orderItems))
		for i, orderItem := range orderItems {
			ItemsReq[i] = &inventory.InventoryReq_Items{
				ProductId: int32(orderItem.ProductId),
				Quantity:  int32(orderItem.Quantity),
			}
		}
		inventoryResp, err := a.InventoryRpc.DecreaseInventory(ctx, &inventory.InventoryReq{
			Items:      ItemsReq,
			PreOrderId: orderModelRes.PreOrderId,
			UserId:     msg.UserID,
		})
		if err != nil {
			logx.Errorw("failed to decrease inventory", logx.Field("err", err), logx.Field("body", string(res.Body)))
			if err := res.Reject(true); err != nil {
				logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
			continue
		}
		// 说明已经扣减过了幂等
		if inventoryResp.StatusCode != code.Success {
			logx.Infow("info to decrease inventory", logx.Field("status_msg", inventoryResp.StatusMsg))
		}
		if orderModelRes.CouponId != "" {
			// 使用优惠券
			couponRes, err := a.CouponRpc.UseCoupon(ctx, &coupons.UseCouponReq{
				CouponId:       orderModelRes.CouponId,
				UserId:         msg.UserID,
				OrderId:        orderModelRes.OrderId,
				PreOrderId:     orderModelRes.PreOrderId,
				DiscountAmount: orderModelRes.DiscountAmount,
				OriginAmount:   orderModelRes.OriginalAmount,
			})
			if err != nil {
				logx.Errorw("failed to use coupon", logx.Field("err", err), logx.Field("body", string(res.Body)))
				if err := res.Reject(true); err != nil {
					logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
				}
				continue
			}
			if couponRes.StatusCode != code.Success {
				logx.Infow("info to use coupon", logx.Field("status_msg", couponRes.StatusMsg))
			}
		}

		if err := res.Ack(false); err != nil {
			logx.Errorw("failed to ack message", logx.Field("err", err), logx.Field("body", string(res.Body)))
		}
	}
}
