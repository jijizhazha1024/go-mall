package product

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductsModel = (*CustomProductsModel)(nil)

type (
	// ProductsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductsModel.
	ProductsModel interface {
		productsModel
		WithSession(session sqlx.Session) ProductsModel
	}

	CustomProductsModel struct {
		*defaultProductsModel
	}
)

// NewProductsModel returns a model for the database table.
func NewProductsModel(conn sqlx.SqlConn) ProductsModel {
	return &CustomProductsModel{
		defaultProductsModel: newProductsModel(conn),
	}
}

func (m *CustomProductsModel) WithSession(session sqlx.Session) ProductsModel {
	return NewProductsModel(sqlx.NewSqlConnFromSession(session))
}
