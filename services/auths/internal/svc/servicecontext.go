package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jijizhazha1024/go-mall/services/auths/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Rdb    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb, err := redis.NewRedis(c.CacheConf[0].RedisConf)
	if err != nil {
		logx.Errorw("redis error", logx.Field("err", err))
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		Rdb:    rdb,
	}
}
