package delay

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
	"time"
)

const (
	ExchangeName = "order-delay-exchange"
	ExchangeKind = "x-delayed-message"
	QueueName    = "order-delay-queue"
	Delay        = 30 * time.Minute
)

type OrderDelayMQ struct {
	conn            *amqp.Connection
	OrderModel      order.OrdersModel
	OrderItemsModel order.OrderItemsModel
	Model           sqlx.SqlConn
	CheckoutRpc     checkoutservice.CheckoutService
	CouponRpc       coupons.CouponsClient
	InventoryRpc    inventory.InventoryClient
}
type OrderReq struct {
	OrderId  string `json:"order_id"`
	UserID   int32  `json:"user_id"`
	RetryCnt int
}

func Init(c config.Config) (*OrderDelayMQ, error) {
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
		ExchangeKind, // 类型为延迟交换机
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部交换机
		false,        // 等待确认
		amqp.Table{
			"x-delayed-type": "direct",
		},
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
	orderDelay := &OrderDelayMQ{
		conn:            conn,
		OrderModel:      order.NewOrdersModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		CheckoutRpc:     checkoutservice.NewCheckoutService(zrpc.MustNewClient(c.CheckoutRpc)),
		CouponRpc:       couponsclient.NewCoupons(zrpc.MustNewClient(c.CouponRpc)),
		InventoryRpc:    inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		Model:           sqlx.NewMysql(c.MysqlConfig.DataSource),
		OrderItemsModel: order.NewOrderItemsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
	}
	go orderDelay.consumer(context.TODO())
	return orderDelay, nil
}
