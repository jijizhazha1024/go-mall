package svc

import (
	"fmt"
	gorse "jijizhazha1024/go-mall/common/utils/gorse"
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/dal/model/user_address"
	"jijizhazha1024/go-mall/services/audit/auditclient"
	"jijizhazha1024/go-mall/services/users/internal/config"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	AuditRpc     auditclient.Audit
	UsersModel   user.UsersModel
	AddressModel user_address.UserAddressesModel
	Model        sqlx.SqlConn
	GorseClient  *gorse.GorseClient
	BF           *bloom.Filter
}

// 初始化监控指标（包级变量改为结构体字段）
var UserRegCounter = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "user_service",
	Subsystem: "register",
	Name:      "total",
	Help:      "Total number of user registrations",
	Labels:    []string{"status"}, // 标签定义
})

func NewServiceContext(c config.Config) *ServiceContext {
	gorseClient := gorse.NewGorseClient(c.GorseConfig.GorseAddr, c.GorseConfig.GorseApikey)
	bf := bloom.New(redis.MustNewRedis(c.RedisConf), "user_login_bloom", 1000000)
	// bloom预热
	usermodel := user.NewUsersModel(sqlx.NewMysql(c.MysqlConfig.DataSource))
	err := bloomPreheat(bf, usermodel)
	if err != nil {
		panic(fmt.Sprintf("bf缓存预热失败: %v", err))
	}

	return &ServiceContext{

		Config:       c,
		GorseClient:  gorseClient,
		Model:        sqlx.NewMysql(c.MysqlConfig.DataSource),
		UsersModel:   usermodel,
		AddressModel: user_address.NewUserAddressesModel(sqlx.NewMysql(c.MysqlConfig.DataSource), c.Cache),
		AuditRpc:     auditclient.NewAudit(zrpc.MustNewClient(c.AuditRpc)),
		BF:           bf,
	}
}
func bloomPreheat(BF *bloom.Filter, UsersModel user.UsersModel) error {

	emails, err := UsersModel.FindAllEmails()
	if err != nil {
		return err
	}

	for _, email := range emails {
		BF.Add([]byte(email))
	}
	return nil

}
