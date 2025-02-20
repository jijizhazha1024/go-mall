package order_items

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderItemsModel = (*customOrderItemsModel)(nil)

type (
	// OrderItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemsModel.
	OrderItemsModel interface {
		orderItemsModel
	}

	customOrderItemsModel struct {
		*defaultOrderItemsModel
	}
)

// NewOrderItemsModel returns a model for the database table.
func NewOrderItemsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderItemsModel {
	return &customOrderItemsModel{
		defaultOrderItemsModel: newOrderItemsModel(conn, c, opts...),
	}
}
