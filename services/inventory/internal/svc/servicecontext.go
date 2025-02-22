package svc

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/dal/model/inventory"
	"jijizhazha1024/go-mall/services/inventory/internal/config"
	"jijizhazha1024/go-mall/services/inventory/internal/mq"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                config.Config
	Rdb                   *redis.Redis
	InventoryModel        inventory.InventoryModel
	InventoryMQ           *mq.InventoryMQ
	DecreaseInventoryShal string
	ReturnInventoryShal   string
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化消息队列等
	inventoryMQ, err := mq.Init(c)
	if err != nil {
		panic(err)
	}

	// 创建ServiceContext实例
	svcCtx := &ServiceContext{
		Config:         c,
		Rdb:            redis.MustNewRedis(c.RedisConf),
		InventoryModel: inventory.NewInventoryModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		InventoryMQ:    inventoryMQ,
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
		err := s.Rdb.Hset(fmt.Sprintf("inventory:%d", inv.ProductId),
			"total", string(rune(inv.Total)))
		if err != nil {
			return fmt.Errorf("缓存库存数据失败: %v", err)
		}
	}

	return nil

}

func (s *ServiceContext) predecreaseloadScript() (string, error) {
	script := `
	    -- 幂等性检查
	    if redis.call("EXISTS", KEYS[1]) == 1 then
	        return 1
	    end

	    -- 预检查库存
	    for i=2, #KEYS do
	        local stock = tonumber(redis.call('GET',KEYS[i]) or 0)
	        local deduct = tonumber(ARGV[i])  -- ARGV索引从1开始
	        if stock < deduct then
			--删除锁
			redis.call("DEL", KEYS[1])
	            return 2
	        end
	    end

	    -- 扣减库存
	    for i=2, #KEYS do
		redis.call('DECRBY', KEYS[i], tonumber(ARGV[i]))
	    end

	    -- 设置处理标记（30分钟过期）
	    redis.call("SET", KEYS[1], ARGV[1], "EX", 1800)
	    return 0
	`
	sha, err := s.Rdb.ScriptLoad(script)

	if err != nil {
		logx.Errorf("Failed to decrease load script: %v", err)
		return "", err
	}
	return sha, nil
}
func (s *ServiceContext) prereturnloadScript() (string, error) {
	script := `
        -- 幂等性检查
        if redis.call("EXISTS", KEYS[1]) == 1 then
            return 1
        end  
        
        -- 归还库存
        for i=2, #KEYS do
            redis.call("INCRBY", KEYS[i], tonumber(ARGV[i]))
        end
        
        -- 设置处理标记（30分钟过期）
        redis.call("SET", KEYS[1], ARGV[1], "EX", 1800)
        return 0
    `

	sha, err := s.Rdb.ScriptLoad(script)

	if err != nil {
		logx.Errorf("Failed to load return script: %v", err)
		return "", err
	}
	return sha, nil
}
