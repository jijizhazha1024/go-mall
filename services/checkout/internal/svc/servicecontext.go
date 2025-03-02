package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/dal/model/cart"
	"jijizhazha1024/go-mall/dal/model/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/config"
	"jijizhazha1024/go-mall/services/checkout/internal/db"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"
	"jijizhazha1024/go-mall/services/inventory/inventoryclient"
	"jijizhazha1024/go-mall/services/product/productcatalogservice"
)

type ServiceContext struct {
	Config             config.Config
	Mysql              sqlx.SqlConn
	RedisClient        *redis.Redis
	CheckoutModel      checkout.CheckoutsModel
	CheckoutItemsModel checkout.CheckoutItemsModel
	CartsModel         cart.CartsModel
	InventoryRpc       inventoryclient.Inventory
	CouponsRpc         couponsclient.Coupons
	ProductRpc         productcatalogservice.ProductCatalogService
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := db.NewMysql(c.MysqlConfig)
	redisconf, _ := redis.NewRedis(c.RedisConf)
	return &ServiceContext{
		Config:             c,
		Mysql:              mysql,
		RedisClient:        redisconf,
		CartsModel:         cart.NewCartsModel(mysql),
		CheckoutModel:      checkout.NewCheckoutsModel(mysql),
		CheckoutItemsModel: checkout.NewCheckoutItemsModel(mysql),
		InventoryRpc:       inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		CouponsRpc:         couponsclient.NewCoupons(zrpc.MustNewClient(c.CouponsRpc)),
		ProductRpc:         productcatalogservice.NewProductCatalogService(zrpc.MustNewClient(c.ProductRpc)),
	}
}
