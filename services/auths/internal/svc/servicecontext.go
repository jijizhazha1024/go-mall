package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/auths/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	UserModel user.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MysqlConfig.DataSource)
	return &ServiceContext{
		UserModel: user.NewUsersModel(conn),
		Config:    c,
	}
}
