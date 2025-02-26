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
		FindProductIsExist(ctx context.Context, productID int64) (bool, error)
		QueryAllProducts(ctx context.Context) ([]*Products, error)
	}

	CustomProductsModel struct {
		*defaultProductsModel
	}
)

func (m *CustomProductsModel) QueryAllProducts(ctx context.Context) ([]*Products, error) {
	query := fmt.Sprintf("SELECT * FROM %s", m.table)
	products := make([]*Products, 0)
	err := m.conn.QueryRowsCtx(ctx, &products, query)
	return products, err
}

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
func (m *defaultProductsModel) FindProductIsExist(ctx context.Context, productID int64) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id=?", m.table)

	err := m.conn.QueryRowCtx(ctx, &count, query, productID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
