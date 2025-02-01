package mq

import (
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/streadway/amqp"
	"time"
)

func (a *AuditMQ) Product(msg *AuditReq) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// 开启事务
	if err := a.channel.Tx(); err != nil {
		return err
	}
	// 发布消息
	// 尝试3次
	if err := retry.Do(func() error {
		return a.channel.Publish(
			ExchangeName,
			RoutingKeyName,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // 持久化消息
				Body:         body,
			},
		)
	}, retry.Attempts(3), retry.Delay(time.Millisecond*100)); err != nil {
		return a.channel.TxRollback()
	}
	// 提交事务
	return a.channel.TxCommit()
}
