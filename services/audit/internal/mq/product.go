package mq

import (
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/streadway/amqp"
	"time"
)

func (a *AuditMQ) Product(msg *AuditReq) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 创建新通道
	channel, err := a.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// 启用事务模式
	if err := channel.Tx(); err != nil {
		return fmt.Errorf("启用事务失败: %v", err)
	}

	// 发布消息到同一通道
	publishErr := retry.Do(
		func() error {
			return channel.Publish(
				ExchangeName,
				RoutingKeyName,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					Body:         body,
				},
			)
		},
		retry.Attempts(3),
		retry.Delay(100*time.Millisecond),
	)

	if publishErr != nil {
		// 回滚事务
		if err := channel.TxRollback(); err != nil {
			return fmt.Errorf("回滚失败: %v, 发布错误: %w", err, publishErr)
		}
		return publishErr
	}

	// 提交事务
	if err := channel.TxCommit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}
