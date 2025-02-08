package mq

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/biz"
	audit2 "jijizhazha1024/go-mall/dal/es/audit"
	"jijizhazha1024/go-mall/dal/model/audit"
	"jijizhazha1024/go-mall/services/audit/internal/config"
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
	channel  *amqp.Channel
	conn     *amqp.Connection
	model    audit.AuditModel
	esClient *elastic.Client
}

func (a *AuditMQ) Close() error {
	if err := a.channel.Close(); err != nil {
		return err
	}
	return a.conn.Close()
}

type AuditReq struct {
	UserID      uint32 `json:"user_id"`
	UserName    string `json:"user_name"`
	ActionType  string `json:"action_type"`
	ActionDesc  string `json:"action_desc"`
	TargetTable string `json:"target_table"`
	TargetID    int64  `json:"target_id"`
	OldData     string `json:"old_data"`
	NewData     string `json:"new_data"`
	// trace
	TraceID   string `json:"trace_id"`
	SpanID    string `json:"span_id"`
	ClientIP  string `json:"client_ip"`
	CreatedAt int64  `json:"created_at"`
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

func initEsIndex(ctx context.Context, client *elastic.Client) error {
	// 创建索引
	exists, err := client.IndexExists(biz.EsIndexName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		createIndex, err := client.CreateIndex(biz.EsIndexName).Body(audit2.Mapping).Do(context.Background())
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			// 处理未确认创建的情况
			return fmt.Errorf("创建索引失败")
		}
	}
	return nil
}
func Init(c config.Config) (*AuditMQ, error) {
	//mysql conn

	model := audit.NewAuditModel(sqlx.NewMysql(c.Mysql.DataSource))
	// es client
	client, err := elastic.NewClient(
		elastic.SetURL(c.ElasticSearch.Addr),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	if err := initEsIndex(context.Background(), client); err != nil {
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
		channel:  channel,
		model:    model,
		esClient: client,
	}
	// 启动监听协程
	if err := mq.consumer(); err != nil {
		return nil, err
	}
	return mq, nil
}
