package svc

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/dal/model/payment"
	"jijizhazha1024/go-mall/services/order/orderservice"
	"jijizhazha1024/go-mall/services/payment/internal/config"
	"log"
)

type ServiceContext struct {
	Config       config.Config
	Rdb          *redis.Redis
	PaymentModel payment.PaymentsModel
	OrderRpc     orderservice.OrderService
	Alipay       *alipay.Client
	PaymentMQ    *amqp.Channel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d%s", c.RabbitMQConfig.User, c.RabbitMQConfig.Pass, c.RabbitMQConfig.Host, c.RabbitMQConfig.Port, c.RabbitMQConfig.VHost))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	// 声明交换机（使用延迟交换机）
	err = ch.ExchangeDeclare(
		"delayed_exchange",  // 交换机名称
		"x-delayed-message", // 类型为延迟交换机
		true,                // 持久化
		false,               // 自动删除
		false,               // 内部交换机
		false,               // 等待确认
		amqp.Table{
			"x-delayed-type": "direct", // 延迟交换机的实际类型
		},
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// 声明队列
	queueName := "delayed_queue"
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 等待确认
		nil,       // 队列参数
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// 绑定队列到交换机
	err = ch.QueueBind(
		queueName,          // 队列名称
		"",                 // 路由键（直接交换机为空）
		"delayed_exchange", // 交换机名称
		false,              // 等待确认
		nil,                // 参数
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}
	// 1. 创建支付宝客户端
	client, _ := alipay.New(c.Alipay.AppId, c.Alipay.PrivateKey, false)
	// 2. 加载支付宝公钥用于验签
	_ = client.LoadAliPayPublicKey(c.Alipay.AlipayPublicKey)

	return &ServiceContext{
		Config:       c,
		Rdb:          redis.MustNewRedis(c.RedisConf),
		PaymentModel: payment.NewPaymentsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		OrderRpc:     orderservice.NewOrderService(zrpc.MustNewClient(c.OrderRpc)),
		Alipay:       client,
		PaymentMQ:    ch,
	}
}
