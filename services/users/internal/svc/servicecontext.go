package svc

import (
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/dal/model/user_address"
	"jijizhazha1024/go-mall/services/users/internal/bloom_filter"
	"jijizhazha1024/go-mall/services/users/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	Bf           *bloom_filter.BloomFilter
	UsersModel   user.UsersModel
	AddressModel user_address.UserAddressesModel
	Model        sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Bf:           bloom_filter.NewBloomFilter(1000000, 0.00001),
		Model:        sqlx.NewMysql(c.MysqlConfig.DataSource),
		UsersModel:   user.NewUsersModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		AddressModel: user_address.NewUserAddressesModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		Config:       c,
	}
}
