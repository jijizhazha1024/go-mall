package svc

import (
	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/dal/model/payment"
	"jijizhazha1024/go-mall/services/order/orderservice"
	"jijizhazha1024/go-mall/services/payment/internal/config"
	"jijizhazha1024/go-mall/services/payment/internal/mq"
)

type ServiceContext struct {
	Config       config.Config
	Rdb          *redis.Redis
	PaymentModel payment.PaymentsModel
	OrderRpc     orderservice.OrderService
	Alipay       *alipay.Client
	PaymentMQ    *mq.PaymentDelayMQ
	Model        sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {

	delayMQ, err := mq.Init(c)
	if err != nil {
		logx.Errorw("创建延迟队列失败", logx.LogField{Key: "err", Value: err})
		panic(err)
	}
	// 1. 创建支付宝客户端
	client, err := alipay.New(c.Alipay.AppId, c.Alipay.PrivateKey, false)
	if err != nil {
		logx.Errorw("创建支付宝客户端失败", logx.LogField{Key: "err", Value: err})
		panic(err)
	}
	// 2. 加载支付宝公钥用于验签
	if err := client.LoadAliPayPublicKey(c.Alipay.AlipayPublicKey); err != nil {
		logx.Errorw("加载支付宝公钥失败", logx.LogField{Key: "err", Value: err})
		panic(err)
	}

	return &ServiceContext{
		Config:       c,
		Rdb:          redis.MustNewRedis(c.RedisConf),
		PaymentModel: payment.NewPaymentsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		OrderRpc:     orderservice.NewOrderService(zrpc.MustNewClient(c.OrderRpc)),
		Alipay:       client,
		PaymentMQ:    delayMQ,
		Model:        sqlx.NewMysql(c.MysqlConfig.DataSource),
	}
}
