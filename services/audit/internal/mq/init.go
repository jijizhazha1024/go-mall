package mq

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"io"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/dal/model/audit"
	"jijizhazha1024/go-mall/services/audit/internal/config"
	"jijizhazha1024/go-mall/services/audit/model/es"
	"strings"
)

const (
	ExchangeName   = "audit_logs:exchange"
	ExchangeType   = amqp.ExchangeDirect
	QueueName      = "audit_logs:queue"
	RoutingKeyName = "audit_logs:routing_key"

	DeadExchangeName   = "audit_logs:dead_exchange"
	DeadQueueName      = "audit_logs:dead_queue"
	DeadRoutingKeyName = "audit_logs:dead_routing_key"
)

type AuditMQ struct {
	conn     *amqp.Connection
	model    audit.AuditModel
	esClient *elasticsearch.Client
}

type AuditReq struct {
	UserID      uint32 `json:"user_id"`
	ActionType  string `json:"action_type"`
	ActionDesc  string `json:"action_desc"`
	TargetTable string `json:"target_table"`
	TargetID    int64  `json:"target_id"`
	OldData     string `json:"old_data"`
	NewData     string `json:"new_data"`
	ServiceName string `json:"service_name"`

	// trace
	TraceID  string `json:"trace_id"`
	SpanID   string `json:"span_id"`
	ClientIP string `json:"client_ip"`

	CreatedAt int64 `json:"created_at"`
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

func initEsIndex(ctx context.Context, client *elasticsearch.Client) error {

	// 创建索引（适配v7客户端）
	existsResp, err := esapi.IndicesExistsRequest{
		Index: []string{biz.EsIndexName},
	}.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("索引存在检查失败: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(existsResp.Body)

	// 检查存在状态
	if existsResp.StatusCode == 404 { // 索引不存在
		createResp, err := esapi.IndicesCreateRequest{
			Index: biz.EsIndexName,
			Body:  strings.NewReader(es.Mapping),
		}.Do(ctx, client)
		if err != nil {
			return fmt.Errorf("创建索引请求失败: %w", err)
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
		}(createResp.Body)
		// 解析响应
		if createResp.IsError() {
			return fmt.Errorf("ES返回错误: %s", createResp.String())
		}
	}
	return nil

}
func Init(c config.Config) (*AuditMQ, error) {
	//mysql conn

	model := audit.NewAuditModel(sqlx.NewMysql(c.MysqlConfig.DataSource))
	// es client

	cfg := elasticsearch.Config{
		Addresses: []string{c.ElasticSearch.Addr},
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	if err := initEsIndex(context.Background(), esClient); err != nil {
		return nil, err
	}
	// mq conn
	conn, err := amqp.Dial(c.RabbitMQ.Dns())
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
	mq := &AuditMQ{
		conn:     conn,
		model:    model,
		esClient: esClient,
	}
	// 启动监听协程
	if err := mq.consumer(); err != nil {
		return nil, err
	}

	return mq, nil
}
