package product_categories

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductCategoriesModel = (*customProductCategoriesModel)(nil)

type (
	// ProductCategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCategoriesModel.
	ProductCategoriesModel interface {
		productCategoriesModel
	}

	customProductCategoriesModel struct {
		*defaultProductCategoriesModel
	}
)

// NewProductCategoriesModel returns a model for the database table.
func NewProductCategoriesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductCategoriesModel {
	return &customProductCategoriesModel{
		defaultProductCategoriesModel: newProductCategoriesModel(conn, c, opts...),
	}
}
