package svc

import (
	_ "embed"
	"fmt"
	"jijizhazha1024/go-mall/dal/model/inventory"
	"jijizhazha1024/go-mall/services/inventory/internal/config"
	"jijizhazha1024/go-mall/services/inventory/internal/decreaselua"
	"jijizhazha1024/go-mall/services/inventory/internal/returnlua"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	Rdb            *redis.Redis
	InventoryModel inventory.InventoryModel
	Locker         rueidislock.Locker

	DecreaseInventoryShal string
	ReturnInventoryShal   string
}

var ProductAccessCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "product_access_total",
		Help: "Total number of product accesses",
	},
	[]string{"product_id", "number"}, // 维度标签
)

func NewServiceContext(c config.Config) *ServiceContext {
	prometheus.MustRegister(ProductAccessCounter)
	lockerOpt := rueidislock.LockerOption{
		ClientOption: rueidis.ClientOption{
			InitAddress: []string{c.RedisConf.Host},
			Password:    c.RedisConf.Pass,
		},
		KeyMajority: 1,               // 单Redis实例模式
		KeyValidity: time.Minute * 5, // 锁有效期
	}

	locker, err := rueidislock.NewLocker(lockerOpt)
	if err != nil {
		panic(fmt.Sprintf("创建Redis锁失败: %v", err))
	}

	// 创建ServiceContext实例
	svcCtx := &ServiceContext{
		Config:         c,
		Rdb:            redis.MustNewRedis(c.RedisConf),
		Locker:         locker,
		InventoryModel: inventory.NewInventoryModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
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
