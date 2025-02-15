package product_categories

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductCategoriesModel = (*customProductCategoriesModel)(nil)

type (
	// ProductCategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCategoriesModel.
	ProductCategoriesModel interface {
		productCategoriesModel
		WithSession(session sqlx.Session) ProductCategoriesModel
		DeleteByProductId(ctx context.Context, productId int64) error
	}

	customProductCategoriesModel struct {
		*defaultProductCategoriesModel
	}
)

// NewProductCategoriesModel returns a model for the database table.
func NewProductCategoriesModel(conn sqlx.SqlConn) ProductCategoriesModel {
	return &customProductCategoriesModel{
		defaultProductCategoriesModel: newProductCategoriesModel(conn),
	}
}

func (m *customProductCategoriesModel) WithSession(session sqlx.Session) ProductCategoriesModel {
	return NewProductCategoriesModel(sqlx.NewSqlConnFromSession(session))
}
func (m *customProductCategoriesModel) DeleteByProductId(ctx context.Context, productId int64) error {
	query := fmt.Sprintf("delete from %s where `product_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, productId)
	return err
}
