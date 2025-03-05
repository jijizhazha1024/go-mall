package notify

import (
	"context"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkoutservice"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/inventory/inventoryclient"
	"jijizhazha1024/go-mall/services/order/internal/config"
)

const (
	ExchangeName = "order-notify-exchange"
	QueueName    = "order-notify-queue"
)

type OrderNotifyMQ struct {
	conn           *amqp.Connection
	OrderModel     order.OrdersModel
	CheckoutRpc    checkoutservice.CheckoutService
	CouponRpc      coupons.CouponsClient
	InventoryRpc   inventory.InventoryClient
	Model          sqlx.SqlConn
	OrderItemModel order.OrderItemsModel
}

type OrderNotifyReq struct {
	OrderId       string `json:"order_id"`
	UserID        int32  `json:"user_id"`
	TransactionId string `json:"transaction_id"`
	PaidAmount    int64  `json:"paid_amount"` // 实际支付金额（分）
	PaidAt        int64  `json:"paid_at"`
}

func Init(c config.Config) (*OrderNotifyMQ, error) {
	conn, err := amqp.Dial(c.RabbitMQConfig.Dns())
	if err != nil {
		return nil, err
	}
	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// 声明交换机（使用延迟交换机）
	err = ch.ExchangeDeclare(
		ExchangeName, // 交换机名称
		amqp.ExchangeDirect,
		true,  // 持久化
		false, // 自动删除
		false, // 内部交换机
		false, // 等待确认
		nil,
	)
	if err != nil {
		return nil, err

	}

	// 声明队列
	_, err = ch.QueueDeclare(
		QueueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 等待确认
		nil,       // 队列参数
	)
	if err != nil {
		return nil, err

	}

	// 绑定队列到交换机
	if err = ch.QueueBind(
		QueueName,
		"",
		ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, err

	}
	orderDelay := &OrderNotifyMQ{
		conn:           conn,
		OrderModel:     order.NewOrdersModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		OrderItemModel: order.NewOrderItemsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		CheckoutRpc:    checkoutservice.NewCheckoutService(zrpc.MustNewClient(c.CheckoutRpc)),
		CouponRpc:      couponsclient.NewCoupons(zrpc.MustNewClient(c.CouponRpc)),
		InventoryRpc:   inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		Model:          sqlx.NewMysql(c.MysqlConfig.DataSource),
	}
	go orderDelay.consumer(context.TODO())
	return orderDelay, nil
}
