package svc

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/services/inventory/inventoryclient"
	"jijizhazha1024/go-mall/services/product/internal/config"
	"net/http"
)

type ServiceContext struct {
	Config       config.Config
	Mysql        sqlx.SqlConn
	RedisClient  *redis.Redis
	Es           *elasticsearch.Client
	InventoryRpc inventoryclient.Inventory
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
		Config:       c,
		Mysql:        mysql,
		RedisClient:  redisconf,
		Es:           es,
		InventoryRpc: inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
	}
}
