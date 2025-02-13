package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/coupons/coupon"
	"jijizhazha1024/go-mall/services/coupons/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	CouponsModel coupon.CouponsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		CouponsModel: coupon.NewCouponsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
	}
}
