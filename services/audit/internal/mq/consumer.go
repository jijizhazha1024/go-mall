package mq

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/dal/model/audit"
	"time"
)

func (a *AuditMQ) consumer() error {
	channel, err := a.conn.Channel()
	if err != nil {
		return err
	}
	results, err := channel.Consume(QueueName, RoutingKeyName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	logx.Infow("consumer start", logx.Field("queue", QueueName), logx.Field("routingKey", RoutingKeyName))
	ctx := context.Background()
	go func() {
		for res := range results {
			var msg *AuditReq
			if err := json.Unmarshal(res.Body, &msg); err != nil {
				logx.Errorw("failed to unmarshal message", logx.Field("error", err), logx.Field("body", string(res.Body)))
				if err := res.Reject(false); err != nil {
					logx.Errorw("failed to reject message", logx.Field("error", err), logx.Field("body", string(res.Body)))
				}
				continue
			}
			// 确保入库
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
		logx.Infow("consumer stopped", logx.Field("queue", QueueName), logx.Field("routingKey", RoutingKeyName))
		defer func() {
			if err := channel.Close(); err != nil {
				logx.Errorw("failed to close channel", logx.Field("error", err))
			}
			if err := a.conn.Close(); err != nil {
				logx.Errorw("failed to close connection", logx.Field("error", err))
			}
		}()
	}()
	return nil
}

func (a *AuditMQ) persistData(ctx context.Context, data *AuditReq) error {
	exist, err := a.model.CheckExistByTraceID(ctx, data.TraceID)
	if err != nil {
		return err
	}
	if !exist {
		if _, err := a.ToMysql(ctx, data); err != nil {
			return err
		}
	}
	if err := a.ToEs(ctx, data); err != nil {
		return err
	}
	return nil
}

func (a *AuditMQ) ToMysql(ctx context.Context, data *AuditReq) (int64, error) {
	res, err := a.model.Insert(ctx, &audit.Audit{
		UserId:      uint64(data.UserID),
		TargetId:    uint64(data.TargetID),
		TargetTable: data.TargetTable,
		ActionType:  data.ActionType,
		ClientIp:    data.ClientIP,
		ActionDesc:  sql.NullString{String: data.ActionDesc, Valid: data.ActionDesc != ""},
		OldData:     sql.NullString{String: data.OldData, Valid: data.OldData != ""},
		NewData:     sql.NullString{String: data.NewData, Valid: data.NewData != ""},
		SpanId:      data.SpanID,
		TraceId:     data.TraceID,
		CreatedAt:   time.Unix(data.CreatedAt, 0),
	})
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *AuditMQ) ToEs(ctx context.Context, data *AuditReq) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	// 创建Index请求
	req := esapi.IndexRequest{
		Index:      biz.EsIndexName,
		DocumentID: data.TraceID,
		Body:       bytes.NewReader(jsonData), // 需要实现AuditReq的JSON序列化方法
	}

	// 执行请求
	res, err := req.Do(ctx, a.esClient)
	if err != nil {
		return fmt.Errorf("IndexRequest failed: %w", err)
	}
	defer res.Body.Close()

	// 检查响应状态
	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}
	return nil
}
