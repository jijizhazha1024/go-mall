package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/payment"
	"jijizhazha1024/go-mall/services/payment/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	Rdb          *redis.Redis
	PaymentModel payment.PaymentsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		Rdb:          redis.MustNewRedis(c.RedisConf),
		PaymentModel: payment.NewPaymentsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
	}
}
