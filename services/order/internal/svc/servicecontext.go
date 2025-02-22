package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkoutservice"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"
	"jijizhazha1024/go-mall/services/order/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	OrderModel  order.OrdersModel
	CheckoutRpc checkoutservice.CheckoutService
	CouponRpc   coupons.CouponsClient
	Model       sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		OrderModel:  order.NewOrdersModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		Model:       sqlx.NewMysql(c.MysqlConfig.DataSource),
		CheckoutRpc: checkoutservice.NewCheckoutService(zrpc.MustNewClient(c.CheckoutRpc)),
		CouponRpc:   couponsclient.NewCoupons(zrpc.MustNewClient(c.CouponRpc)),
	}
}
