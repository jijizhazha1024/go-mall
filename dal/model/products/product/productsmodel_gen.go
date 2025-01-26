// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.5

package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	productsFieldNames          = builder.RawFieldNames(&Products{})
	productsRows                = strings.Join(productsFieldNames, ",")
	productsRowsExpectAutoSet   = strings.Join(stringx.Remove(productsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	productsRowsWithPlaceHolder = strings.Join(stringx.Remove(productsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheProductsIdPrefix = "cache:products:id:"
)

type (
	productsModel interface {
		Insert(ctx context.Context, data *Products) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Products, error)
		Update(ctx context.Context, data *Products) error
		Delete(ctx context.Context, id int64) error
	}

	defaultProductsModel struct {
		sqlc.CachedConn
		table string
	}

	Products struct {
		Id          int64          `db:"id"`          // 主键，自增,商品id
		NAME        string         `db:"NAME"`        // 商品名称
		Description sql.NullString `db:"description"` // 商品描述
		Picture     sql.NullString `db:"picture"`     // 商品图片信息
		Price       float64        `db:"price"`       // 商品价格
		Stock       int64          `db:"stock"`       // 库存数量
		CreatedAt   time.Time      `db:"created_at"`  // 创建时间
		UpdatedAt   time.Time      `db:"updated_at"`  // 更新时间
	}
)

func newProductsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultProductsModel {
	return &defaultProductsModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`products`",
	}
}

func (m *defaultProductsModel) Delete(ctx context.Context, id int64) error {
	productsIdKey := fmt.Sprintf("%s%v", cacheProductsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, productsIdKey)
	return err
}

func (m *defaultProductsModel) FindOne(ctx context.Context, id int64) (*Products, error) {
	productsIdKey := fmt.Sprintf("%s%v", cacheProductsIdPrefix, id)
	var resp Products
	err := m.QueryRowCtx(ctx, &resp, productsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productsRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProductsModel) Insert(ctx context.Context, data *Products) (sql.Result, error) {
	productsIdKey := fmt.Sprintf("%s%v", cacheProductsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, productsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.NAME, data.Description, data.Picture, data.Price, data.Stock)
	}, productsIdKey)
	return ret, err
}

func (m *defaultProductsModel) Update(ctx context.Context, data *Products) error {
	productsIdKey := fmt.Sprintf("%s%v", cacheProductsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, productsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.NAME, data.Description, data.Picture, data.Price, data.Stock, data.Id)
	}, productsIdKey)
	return err
}

func (m *defaultProductsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheProductsIdPrefix, primary)
}

func (m *defaultProductsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultProductsModel) tableName() string {
	return m.table
}
