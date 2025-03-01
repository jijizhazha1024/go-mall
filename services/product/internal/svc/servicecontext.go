package svc

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/common/consts/biz"
	gorse "jijizhazha1024/go-mall/common/utils/gorse"
	"jijizhazha1024/go-mall/dal/es/product"
	"jijizhazha1024/go-mall/dal/model/products/categories"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
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
	GorseClient     *gorse.GorseClient
	ProductModel    product2.ProductsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 Redis 配置
	redisClient, err := redis.NewRedis(c.RedisConf)
	if err != nil {
		logx.Errorw("redis init error", logx.Field("err", err))
		panic(err)
	}
	// 初始化 ES 客户端
	client, err := elastic.NewClient(elastic.SetURL(c.ElasticSearch.Addr),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeoutStartup(30*time.Second))
	if err != nil {
		logx.Errorw("elasticsearch init error", logx.Field("err", err))
		panic(err)
	}
	if err := initEs(context.TODO(), client); err != nil {
		logx.Errorw("elasticsearch init index error", logx.Field("err", err))
		panic(err)
	}
	gorseClient := gorse.NewGorseClient(c.GorseConfig.GorseAddr, c.GorseConfig.GorseApikey)
	return &ServiceContext{
		Config:          c,
		Mysql:           sqlx.NewMysql(c.MysqlConfig.DataSource),
		RedisClient:     redisClient,
		EsClient:        client,
		GorseClient:     gorseClient,
		ProductModel:    product2.NewProductsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		InventoryRpc:    inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		CategoriesModel: categories.NewCategoriesModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
	}
}
func initEs(ctx context.Context, esClient *elastic.Client) error {
	exists, err := esClient.IndexExists(biz.ProductEsIndexName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		createIndex, err := esClient.CreateIndex(biz.ProductEsIndexName).Body(product.EsMapping).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return err
		}
	}
	return nil
}
