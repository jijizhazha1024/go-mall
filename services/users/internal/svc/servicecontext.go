package svc

import (
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/dal/model/user_address"
	"jijizhazha1024/go-mall/services/users/internal/bloom_filter"
	"jijizhazha1024/go-mall/services/users/internal/config"
	"jijizhazha1024/go-mall/services/users/internal/db"
)

type ServiceContext struct {
	Config       config.Config
	Bf           *bloom_filter.BloomFilter
	UsersModel   user.UsersModel
	AddressModel user_address.UserAddressesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := db.NewMysql(c.MysqlConfig)

	return &ServiceContext{
		Config:       c,
		AddressModel: user_address.NewUserAddressesModel(mysql),
		UsersModel:   user.NewUsersModel(mysql),
		Bf:           bloom_filter.NewBloomFilter(1000000, 0.00001),
	}

}
