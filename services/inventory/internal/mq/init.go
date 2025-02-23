package mq

import (
	"github.com/streadway/amqp"
	"jijizhazha1024/go-mall/services/inventory/internal/config"
)

const (
	ExchangeName   = "inventory:exchange"
	ExchangeType   = amqp.ExchangeDirect
	QueueName      = "inventory:queue"
	RoutingKeyName = "inventory:routing_key"

	DeadExchangeName   = "inventory:dead_exchange"
	DeadQueueName      = "inventory:dead_queue"
	DeadRoutingKeyName = "inventory:dead_routing_key"
)

type InventoryMQ struct {
	conn *amqp.Connection
}

type InventoryReq struct {
	Quantity  int64
	ProductId int64
}

func declareMainQueue(channel *amqp.Channel) error {
	if err := channel.ExchangeDeclare(ExchangeName, ExchangeType,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	); err != nil {
		return err
	}

	// 声明队列（带死信配置）
	if _, err := channel.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    DeadExchangeName,
			"x-dead-letter-routing-key": DeadRoutingKeyName,
		},
	); err != nil {
		return err
	}
	// 绑定主队列
	if err := channel.QueueBind(QueueName, RoutingKeyName, ExchangeName, false, nil); err != nil {
		return err
	}
	return nil
}
func declareDeadQueue(channel *amqp.Channel) error {
	// 声明死信交换机
	if err := channel.ExchangeDeclare(DeadExchangeName, amqp.ExchangeDirect, true, false, false, false,
		nil); err != nil {
		return err
	}

	// 声明死信队列
	if _, err := channel.QueueDeclare(DeadQueueName, true, false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	// 绑定死信队列
	if err := channel.QueueBind(
		DeadQueueName,
		DeadRoutingKeyName,
		DeadExchangeName,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func Init(c config.Config) (*InventoryMQ, error) {
	//mysql conn

	// mq conn
	conn, err := amqp.Dial(c.RabbitMQConfig.Dns())
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func(channel *amqp.Channel) {
		err = channel.Close()
	}(channel)
	// 声明交换机
	if err := declareMainQueue(channel); err != nil {
		return nil, err
	}
	if err := declareDeadQueue(channel); err != nil {
		return nil, err
	}
	// 启动监听协程
	mq := &InventoryMQ{
		conn: conn,
	}
	// 启动监听协程
	if err := mq.consumer(); err != nil {
		return nil, err
	}

	return mq, nil
}
