package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/order/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel order.OrdersModel
	Model      sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		OrderModel: order.NewOrdersModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		Model:      sqlx.NewMysql(c.MysqlConfig.DataSource),
	}
}
