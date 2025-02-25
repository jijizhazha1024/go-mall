package svc

import (
	"fmt"
	"github.com/olivere/elastic"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/dal/model/products/categories"
	"jijizhazha1024/go-mall/services/inventory/inventoryclient"
	"jijizhazha1024/go-mall/services/product/internal/config"
	"time"
)

type ServiceContext struct {
	Config          config.Config
	Mysql           sqlx.SqlConn
	RedisClient     *redis.Redis
	CategoriesModel categories.CategoriesModel
	EsClient        *elastic.Client
	InventoryRpc    inventoryclient.Inventory
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := sqlx.NewMysql(c.MysqlConfig.DataSource)
	// 初始化 Redis 配置
	redisconf, _ := redis.NewRedis(c.RedisConf)
	// 初始化 ES 客户端
	fmt.Println(c.ElasticSearch.Addr)
	client, err := elastic.NewClient(elastic.SetURL(c.ElasticSearch.Addr),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeoutStartup(30*time.Second))
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:          c,
		Mysql:           mysql,
		RedisClient:     redisconf,
		EsClient:        client,
		InventoryRpc:    inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		CategoriesModel: categories.NewCategoriesModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
	}
}
