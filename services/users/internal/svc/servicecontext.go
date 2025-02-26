package svc

import (
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/dal/model/user_address"
	"jijizhazha1024/go-mall/services/audit/audit"
	"jijizhazha1024/go-mall/services/audit/auditclient"
	"jijizhazha1024/go-mall/services/users/internal/config"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	AuditRpc     audit.AuditClient
	UsersModel   user.UsersModel
	AddressModel user_address.UserAddressesModel
	Model        sqlx.SqlConn
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

	return &ServiceContext{

		Model:        sqlx.NewMysql(c.MysqlConfig.DataSource),
		UsersModel:   user.NewUsersModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		AddressModel: user_address.NewUserAddressesModel(sqlx.NewMysql(c.MysqlConfig.DataSource), c.Cache),
		AuditRpc:     auditclient.NewAudit(zrpc.MustNewClient(c.AuditRpc)),

		Config: c,
	}
}
