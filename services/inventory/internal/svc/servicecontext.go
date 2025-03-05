package svc

import (
	"context"
	_ "embed"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/dal/model/inventory"
	"jijizhazha1024/go-mall/services/inventory/internal/config"
	"jijizhazha1024/go-mall/services/inventory/internal/decreaselua"
	"jijizhazha1024/go-mall/services/inventory/internal/returnlua"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	Rdb            *redis.Redis
	InventoryModel inventory.InventoryModel

	DecreaseInventoryShal string
	ReturnInventoryShal   string
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
	decreaseInventoryShashal, err := svcCtx.predecreaseloadScript()
	if err != nil {
		panic(fmt.Sprintf("加载Lua脚本失败: %v", err))
	}
	svcCtx.DecreaseInventoryShal = decreaseInventoryShashal
	returnInventoryShashal, err := svcCtx.prereturnloadScript()
	if err != nil {
		panic(fmt.Sprintf("加载Lua脚本失败: %v", err))
	}
	svcCtx.ReturnInventoryShal = returnInventoryShashal

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
		productKey := fmt.Sprintf("%s:%d", biz.InventoryProductKey, inv.ProductId)
		if err := s.Rdb.Set(productKey, strconv.Itoa(int(inv.Total))); err != nil {
			return fmt.Errorf("缓存库存数据失败: %v", err)
		}
	}
	return nil

}

func (s *ServiceContext) predecreaseloadScript() (string, error) {

	sha, err := s.Rdb.ScriptLoad(decreaselua.Decreaselua)

	if err != nil {
		logx.Errorf("Failed to decrease load script: %v", err)
		return "", err
	}
	return sha, nil
}
func (s *ServiceContext) prereturnloadScript() (string, error) {

	sha, err := s.Rdb.ScriptLoad(returnlua.Returnlua)

	if err != nil {
		logx.Errorf("Failed to load return script: %v", err)
		return "", err
	}
	return sha, nil
}
