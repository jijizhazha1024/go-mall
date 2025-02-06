package svc

import (
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/bloom_filter"
	"jijizhazha1024/go-mall/services/users/internal/config"
	"jijizhazha1024/go-mall/services/users/internal/db"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	Mysql      sqlx.SqlConn
	Bf         *bloom_filter.BloomFilter
	UsersModel user.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := db.NewMysql(c.MysqlConfig)
	return &ServiceContext{
		Config:     c,
		Mysql:      mysql,
		UsersModel: user.NewUsersModel(mysql),
	}
}
