package categories

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CategoriesModel = (*customCategoriesModel)(nil)

type (
	// CategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoriesModel.
	CategoriesModel interface {
		categoriesModel
		withSession(session sqlx.Session) CategoriesModel
		FindCategoryNameByProductID(ctx context.Context, productID int64) ([]string, error)
	}

	customCategoriesModel struct {
		*defaultCategoriesModel
	}
)

// NewCategoriesModel returns a model for the database table.
func NewCategoriesModel(conn sqlx.SqlConn) CategoriesModel {
	return &customCategoriesModel{
		defaultCategoriesModel: newCategoriesModel(conn),
	}
}

func (m *customCategoriesModel) withSession(session sqlx.Session) CategoriesModel {
	return NewCategoriesModel(sqlx.NewSqlConnFromSession(session))
}

// FindByProductID 根据商品ID查询商品分类
func (m *defaultCategoriesModel) FindCategoryNameByProductID(ctx context.Context, productID int64) ([]string, error) {
	query := fmt.Sprintf(
		"SELECT c.name FROM %s as c INNER JOIN %s pc ON c.id = pc.category_id WHERE pc.product_id = ?",
		m.table, "product_categories")
	var categories []string
	err := m.conn.QueryRowsCtx(ctx, &categories, query, productID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
