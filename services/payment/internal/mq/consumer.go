package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

func (a *PaymentDelayMQ) consumer(ctx context.Context) {

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
		var msg *PaymentReq
		if err := json.Unmarshal(res.Body, &msg); err != nil {
			logx.Errorw("failed to unmarshal message", logx.Field("error", err), logx.Field("body", string(res.Body)))
			if err := res.Reject(false); err != nil {
				logx.Errorw("failed to reject message", logx.Field("error", err), logx.Field("body", string(res.Body)))
			}
			continue
		}
		fmt.Println(msg)

	}
}
