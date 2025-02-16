package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/cart"
	"jijizhazha1024/go-mall/services/carts/internal/config"
	"jijizhazha1024/go-mall/services/carts/internal/db"
)

type ServiceContext struct {
	Config     config.Config
	Mysql      sqlx.SqlConn
	CartsModel cart.CartsModel
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	mysql := db.NewMysql(c.MysqlConfig)
	return &ServiceContext{
		Config:     c,
		Mysql:      mysql,
		CartsModel: cart.NewCartsModel(mysql),
	}, nil
}
