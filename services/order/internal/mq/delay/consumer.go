package delay

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/order/order"
)

func (a *OrderDelayMQ) consumer(ctx context.Context) {
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

		var msg *OrderReq
		var orderModelRes *order2.Orders
		if err := json.Unmarshal(res.Body, &msg); err != nil {
			logx.Errorw("failed to unmarshal message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			if err := res.Reject(true); err != nil {
				logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
			continue
		}
		// --------------- reverse --------------- 幂等
		// 1. 更新订单状态为已过期
		// 2. 释放优惠券
		// 3. 释放预扣减的库存
		isContinue := true
		if err := a.Model.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
			ordersModel := a.OrderModel.WithSession(session)
			orderRes, err := ordersModel.GetOrderByOrderIDAndUserIDWithLock(ctx, msg.OrderId, msg.UserID)
			if err != nil {
				return err
			}
			orderModelRes = orderRes
			
			// 只进行处理创建订单的订单
			if order.OrderStatus(orderRes.OrderStatus) != order.OrderStatus_ORDER_STATUS_CREATED {
				isContinue = false
				return nil
			}
			if err := ordersModel.UpdateOrderStatusByOrderIDAndUserID(
				ctx,
				msg.OrderId,
				msg.UserID,
				order.OrderStatus_ORDER_STATUS_CLOSED,
				order.PaymentStatus_PAYMENT_STATUS_EXPIRED,
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			logx.Errorw("failed to update order status", logx.Field("err", err), logx.Field("body", string(res.Body)))
			if err := res.Reject(true); err != nil {
				logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
		}
		if !isContinue {
			logx.Infow("info to update order status")
			if err := res.Ack(true); err != nil {
				logx.Errorw("failed to ack message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
			continue
		}
		orderItems, err := a.OrderItemsModel.QueryOrderItemsByOrderID(ctx, orderModelRes.OrderId)
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
		returnPreInventoryResp, err := a.InventoryRpc.ReturnPreInventory(ctx, &inventory.InventoryReq{
			PreOrderId: orderModelRes.PreOrderId,
			Items:      ItemsReq,
			UserId:     int32(orderModelRes.UserId),
		})
		if err != nil {
			logx.Errorw("failed to decrease pre inventory", logx.Field("err", err), logx.Field("body", string(res.Body)))
			if err := res.Reject(true); err != nil {
				logx.Errorw("failed to reject message", logx.Field("err", err), logx.Field("body", string(res.Body)))
			}
			continue
		}
		if returnPreInventoryResp.StatusCode != code.Success {
			logx.Infow("info to decrease pre inventory", logx.Field("status_msg", returnPreInventoryResp.StatusMsg))
		}
		if err := res.Ack(false); err != nil {
			logx.Errorw("failed to ack message", logx.Field("err", err), logx.Field("body", string(res.Body)))
		}
		logx.Infow("consumer order", logx.Field("body", string(res.Body)),
			logx.Field("queue", QueueName))
	}
}
