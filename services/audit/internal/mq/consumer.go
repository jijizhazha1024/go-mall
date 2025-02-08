package mq

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/dal/model/audit"
	"time"
)

func (a *AuditMQ) consumer() error {
	results, err := a.channel.Consume(QueueName, RoutingKeyName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	logx.Infow("consumer start", logx.Field("queue", QueueName), logx.Field("routingKey", RoutingKeyName))
	ctx := context.Background()
	go func() {
		for res := range results {
			var msg *AuditReq
			if err := json.Unmarshal(res.Body, &msg); err != nil {
				logx.Errorw("failed to unmarshal message", logx.Field("error", err))
				continue
			}
			// 确保入库
			// TODO：可能存在重复消息问题，需要优化
			if err := a.persistData(ctx, msg); err != nil {
				logx.Errorw("insert failed, rejecting message", logx.Field("err", err), logx.Field("msg", msg))
				// 显式拒绝消息，不重新入队（requeue=false），进入死信队列
				if err := res.Nack(false, false); err != nil {
					logx.Errorw("NACK failed", logx.Field("err", err), logx.Field("msg", msg))
				}
				continue
			}
			// 消息确认
			if err := res.Ack(false); err != nil {
				logx.Errorw("ACK failed", logx.Field("error", err), logx.Field("body", string(res.Body))) // 记录完整消息内容
			}
		}
	}()
	return nil
}

func (a *AuditMQ) persistData(ctx context.Context, data *AuditReq) error {
	if err := a.ToMysql(ctx, data); err != nil {
		return err
	}
	if err := a.ToEs(ctx, data); err != nil {
		return err
	}
	return nil
}

func (a *AuditMQ) ToMysql(ctx context.Context, data *AuditReq) error {
	if _, err := a.model.Insert(ctx, &audit.Audit{
		UserId:      uint64(data.UserID),
		Username:    data.UserName,
		TargetId:    uint64(data.TargetID),
		TargetTable: data.TargetTable,
		ActionType:  data.ActionType,
		ClientIp:    sql.NullString{String: data.ClientIP, Valid: true},
		ActionDesc:  sql.NullString{String: data.ActionDesc, Valid: true},
		OldData:     sql.NullString{String: data.OldData, Valid: true},
		NewData:     sql.NullString{String: data.NewData, Valid: true},
		SpanId:      sql.NullString{String: data.SpanID, Valid: true},
		TraceId:     sql.NullString{String: data.TraceID, Valid: true},
		CreatedAt:   time.Unix(data.CreatedAt, 0),
	}); err != nil {
		return err
	}
	return nil
}
func (a *AuditMQ) ToEs(ctx context.Context, data *AuditReq) error {
	// 使用Elasticsearch客户端的Index()方法插入文档
	if _, err := a.esClient.Index().
		Index(biz.EsIndexName).
		BodyJson(*data). // 或者使用BodyString()如果你已经有了JSON字符串
		Do(ctx); err != nil {
		return err
	}
	return nil
}
