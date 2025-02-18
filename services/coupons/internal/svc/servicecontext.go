package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/dal/model/coupons/coupon"
	"jijizhazha1024/go-mall/dal/model/coupons/coupon_usage"
	"jijizhazha1024/go-mall/dal/model/coupons/user_coupons"
	"jijizhazha1024/go-mall/services/coupons/internal/config"
)

type ServiceContext struct {
	Config           config.Config
	CouponsModel     coupon.CouponsModel
	UserCouponsModel user_coupons.UserCouponsModel
	CouponUsageModel coupon_usage.CouponUsageModel
	Model            sqlx.SqlConn
	Rdb              *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		CouponsModel:     coupon.NewCouponsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		UserCouponsModel: user_coupons.NewUserCouponsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		CouponUsageModel: coupon_usage.NewCouponUsageModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		Model:            sqlx.NewMysql(c.MysqlConfig.DataSource),
		Rdb:              redis.MustNewRedis(c.RedisConf),
	}
}
