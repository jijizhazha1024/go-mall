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
		FindCategoriesByIds(ctx context.Context, productId int64) ([]string, error)
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
func (m *customProductCategoriesModel) FindCategoriesByIds(ctx context.Context, productId int64) ([]string, error) {
	query := fmt.Sprintf("select `category_id` from %s where `product_id` = ?", m.table)

	// 定义一个切片用于存储查询结果
	var categoryIds []string

	// 使用 QueryRowsCtx 执行查询
	err := m.conn.QueryRowsCtx(ctx, &categoryIds, query, productId)
	if err != nil {
		return nil, err // 如果查询失败，返回错误
	}

	return categoryIds, nil // 返回分类 ID 切片
}
