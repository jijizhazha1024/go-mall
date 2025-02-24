package svc

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/dal/model/inventory"
	"jijizhazha1024/go-mall/services/inventory/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	Rdb            *redis.Redis
	InventoryModel inventory.InventoryModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	// 创建ServiceContext实例
	svcCtx := &ServiceContext{
		Config:         c,
		Rdb:            redis.MustNewRedis(c.RedisConf),
		InventoryModel: inventory.NewInventoryModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
	}

	// 执行缓存预热
	if err := svcCtx.PreheatInventoryCache(); err != nil {
		panic(fmt.Sprintf("缓存预热失败: %v", err))
	}

	return svcCtx
}

// 新增预热方法
func (s *ServiceContext) PreheatInventoryCache() error {
	// 1. 从数据库读取所有库存数据（或指定商品）
	inventories, err := s.InventoryModel.FindAll(context.Background())
	if err != nil {
		return fmt.Errorf("读取库存数据失败: %v", err)
	}
	// 2. 缓存库存数据

	for _, inv := range inventories {
		err := s.Rdb.Hset(fmt.Sprintf("inventory:%d", inv.ProductId),
			"total", string(rune(inv.Total)))
		if err != nil {
			return fmt.Errorf("缓存库存数据失败: %v", err)
		}
	}

	return nil

}
