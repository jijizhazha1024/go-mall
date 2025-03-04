package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkoutservice"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/coupons/couponsclient"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/inventory/inventoryclient"
	"jijizhazha1024/go-mall/services/order/internal/config"
	"jijizhazha1024/go-mall/services/order/internal/mq/delay"
	"jijizhazha1024/go-mall/services/order/internal/mq/notify"
	"jijizhazha1024/go-mall/services/users/users"
	"jijizhazha1024/go-mall/services/users/usersclient"
)

type ServiceContext struct {
	Config         config.Config
	OrderModel     order.OrdersModel
	OrderItemModel order.OrderItemsModel
	OrderAddress   order.OrderAddressesModel
	CheckoutRpc    checkoutservice.CheckoutService
	CouponRpc      coupons.CouponsClient
	UserRpc        users.UsersClient
	InventoryRpc   inventory.InventoryClient
	Model          sqlx.SqlConn
	OrderDelayMQ   *delay.OrderDelayMQ
	OrderNotifyMQ  *notify.OrderNotifyMQ
}

func NewServiceContext(c config.Config) *ServiceContext {
	orderDelayMQ, err := delay.Init(c)
	if err != nil {
		logx.Error(err)
		panic(err)
	}
	notifyMQ, err := notify.Init(c)
	if err != nil {
		logx.Error(err)
		panic(err)
	}
	return &ServiceContext{
		Config:         c,
		OrderModel:     order.NewOrdersModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		OrderItemModel: order.NewOrderItemsModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		OrderAddress:   order.NewOrderAddressesModel(sqlx.NewMysql(c.MysqlConfig.DataSource)),
		Model:          sqlx.NewMysql(c.MysqlConfig.DataSource),
		CheckoutRpc:    checkoutservice.NewCheckoutService(zrpc.MustNewClient(c.CheckoutRpc)),
		CouponRpc:      couponsclient.NewCoupons(zrpc.MustNewClient(c.CouponRpc)),
		UserRpc:        usersclient.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		InventoryRpc:   inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		OrderDelayMQ:   orderDelayMQ,
		OrderNotifyMQ:  notifyMQ,
	}
}
