package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/inventory"
	"jijizhazha1024/go-mall/services/inventory/internal/config"
	"jijizhazha1024/go-mall/services/inventory/internal/mq"
)

type ServiceContext struct {
	Config         config.Config
	Rdb            *redis.Redis
	InventoryModel inventory.InventoryModel
	InventoryMQ    *mq.InventoryMQ
}

func NewServiceContext(c config.Config) *ServiceContext {
	inventoryMQ, err := mq.Init(c)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:         c,
		Rdb:            redis.MustNewRedis(c.RedisConf),
		InventoryModel: inventory.NewInventoryModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		InventoryMQ:    inventoryMQ,
	}
}
