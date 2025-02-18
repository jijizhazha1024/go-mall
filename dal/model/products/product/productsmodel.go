package product

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductsModel = (*CustomProductsModel)(nil)

type (
	// ProductsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductsModel.
	ProductsModel interface {
		productsModel
		WithSession(session sqlx.Session) ProductsModel
		FindPage(ctx context.Context, offset, limit int) ([]*Products, error)
		Count(ctx context.Context) (int64, error)
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

func (m *defaultProductsModel) FindPage(ctx context.Context, offset, limit int) ([]*Products, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", m.table)
	var products []*Products
	err := m.conn.QueryRowsCtx(ctx, &products, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (m *defaultProductsModel) Count(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
