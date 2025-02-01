package mq

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/dal/model/audit"
	"time"
)

func (a *AuditMQ) consumer() error {
	results, err := a.channel.Consume(QueueName, RoutingKeyName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	logx.Infow("consumer start", logx.Field("queue", QueueName), logx.Field("routingKey", RoutingKeyName))
	go func() {
		for res := range results {
			var msg *AuditReq
			if err := json.Unmarshal(res.Body, &msg); err != nil {
				logx.Errorw("failed to unmarshal message", logx.Field("error", err))
				continue
			}
			if _, err := a.model.Insert(context.Background(), &audit.Audit{
				UserId:      uint64(msg.UserID),
				Username:    msg.UserName,
				TargetId:    uint64(msg.TargetID),
				TargetTable: msg.TargetTable,
				ActionType:  msg.ActionType,
				ClientIp:    sql.NullString{String: msg.ClientIP, Valid: true},
				ActionDesc:  sql.NullString{String: msg.ActionDesc, Valid: true},
				OldData:     sql.NullString{String: msg.OldData, Valid: true},
				NewData:     sql.NullString{String: msg.NewData, Valid: true},
				SpanId:      sql.NullString{String: msg.SpanID, Valid: true},
				TraceId:     sql.NullString{String: msg.TraceID, Valid: true},
				CreatedAt:   time.Unix(msg.CreatedAt, 0),
			}); err != nil {
				logx.Errorw("insert failed, rejecting message",
					logx.Field("err", err),
					logx.Field("body", string(res.Body)),
				)
				// 显式拒绝消息，不重新入队（requeue=false），进入死信队列
				if err := res.Nack(false, false); err != nil {
					logx.Errorw("NACK failed",
						logx.Field("err", err),
						logx.Field("body", string(res.Body)),
					)
				}
				continue
			}
			// 消息确认
			if err := res.Ack(false); err != nil {
				logx.Errorw("ACK failed",
					logx.Field("error", err),
					logx.Field("body", string(res.Body)), // 记录完整消息内容
				)
			}
		}
	}()
	return nil
}
