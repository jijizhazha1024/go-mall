package svc

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/services/product/internal/config"
	"net/http"
)

type ServiceContext struct {
	Config      config.Config
	Mysql       sqlx.SqlConn
	RedisClient *redis.Redis
	Es          *elasticsearch.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := sqlx.NewMysql(c.MysqlConfig.DataSource)
	// 初始化 Redis 配置
	redisconf, _ := redis.NewRedis(c.RedisConf)
	// 初始化 ES 客户端
	es, _ := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: c.ElasticsearchConfig.Addresses,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   c.ElasticsearchConfig.MaxIdleConnsPerHost,
			ResponseHeaderTimeout: c.ElasticsearchConfig.ResponseHeaderTimeout,
		},
	})
	return &ServiceContext{
		Config:      c,
		RedisClient: redisconf,
		Mysql:       mysql,
		Es:          es,
	}
}
