package order_addresses

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderAddressesModel = (*customOrderAddressesModel)(nil)

type (
	// OrderAddressesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderAddressesModel.
	OrderAddressesModel interface {
		orderAddressesModel
	}

	customOrderAddressesModel struct {
		*defaultOrderAddressesModel
	}
)

// NewOrderAddressesModel returns a model for the database table.
func NewOrderAddressesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderAddressesModel {
	return &customOrderAddressesModel{
		defaultOrderAddressesModel: newOrderAddressesModel(conn, c, opts...),
	}
}
