// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.5

package product_categories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	productCategoriesFieldNames          = builder.RawFieldNames(&ProductCategories{})
	productCategoriesRows                = strings.Join(productCategoriesFieldNames, ",")
	productCategoriesRowsExpectAutoSet   = strings.Join(stringx.Remove(productCategoriesFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	productCategoriesRowsWithPlaceHolder = strings.Join(stringx.Remove(productCategoriesFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	productCategoriesModel interface {
		Insert(ctx context.Context, data *ProductCategories) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*ProductCategories, error)
		FindOneByProductIdCategoryId(ctx context.Context, productId sql.NullInt64, categoryId sql.NullInt64) (*ProductCategories, error)
		Update(ctx context.Context, data *ProductCategories) error
		Delete(ctx context.Context, id int64) error
	}

	defaultProductCategoriesModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ProductCategories struct {
		Id         int64         `db:"id"`          // 自增主键
		ProductId  sql.NullInt64 `db:"product_id"`  // 商品id
		CategoryId sql.NullInt64 `db:"category_id"` // 分类id
	}
)

func newProductCategoriesModel(conn sqlx.SqlConn) *defaultProductCategoriesModel {
	return &defaultProductCategoriesModel{
		conn:  conn,
		table: "`product_categories`",
	}
}

func (m *defaultProductCategoriesModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultProductCategoriesModel) FindOne(ctx context.Context, id int64) (*ProductCategories, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productCategoriesRows, m.table)
	var resp ProductCategories
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProductCategoriesModel) FindOneByProductIdCategoryId(ctx context.Context, productId sql.NullInt64, categoryId sql.NullInt64) (*ProductCategories, error) {
	var resp ProductCategories
	query := fmt.Sprintf("select %s from %s where `product_id` = ? and `category_id` = ? limit 1", productCategoriesRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, productId, categoryId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProductCategoriesModel) Insert(ctx context.Context, data *ProductCategories) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, productCategoriesRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ProductId, data.CategoryId)
	return ret, err
}

func (m *defaultProductCategoriesModel) Update(ctx context.Context, newData *ProductCategories) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, productCategoriesRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.ProductId, newData.CategoryId, newData.Id)
	return err
}

func (m *defaultProductCategoriesModel) tableName() string {
	return m.table
}
