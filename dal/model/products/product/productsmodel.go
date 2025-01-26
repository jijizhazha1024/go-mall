package product

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductsModel = (*customProductsModel)(nil)

type (
	// ProductsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductsModel.
	ProductsModel interface {
		productsModel
	}

	customProductsModel struct {
		*defaultProductsModel
	}
)

// NewProductsModel returns a model for the database table.
func NewProductsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductsModel {
	return &customProductsModel{
		defaultProductsModel: newProductsModel(conn, c, opts...),
	}
}
