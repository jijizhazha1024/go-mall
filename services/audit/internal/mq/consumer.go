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
				res.Reject(false)
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
				res.Reject(false)
				continue
			}
			// 消息确认
			res.Ack(false)
		}
	}()
	return nil
}
