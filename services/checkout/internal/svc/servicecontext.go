package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/config"
	"jijizhazha1024/go-mall/services/checkout/internal/db"
)

type ServiceContext struct {
	Config        config.Config
	Mysql         sqlx.SqlConn
	RedisClient   *redis.Redis
	CheckoutModel checkout.CheckoutsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := db.NewMysql(c.MysqlConfig)
	redisconf, _ := redis.NewRedis(c.RedisConf)
	return &ServiceContext{
		Config:        c,
		Mysql:         mysql,
		RedisClient:   redisconf,
		CheckoutModel: checkout.NewCheckoutsModel(mysql),
	}
}
