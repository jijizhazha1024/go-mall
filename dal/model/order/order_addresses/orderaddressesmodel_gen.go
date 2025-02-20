// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.5

package order_addresses

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
	orderAddressesFieldNames          = builder.RawFieldNames(&OrderAddresses{})
	orderAddressesRows                = strings.Join(orderAddressesFieldNames, ",")
	orderAddressesRowsExpectAutoSet   = strings.Join(stringx.Remove(orderAddressesFieldNames, "`address_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	orderAddressesRowsWithPlaceHolder = strings.Join(stringx.Remove(orderAddressesFieldNames, "`address_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheOrderAddressesAddressIdPrefix = "cache:orderAddresses:addressId:"
)

type (
	orderAddressesModel interface {
		Insert(ctx context.Context, data *OrderAddresses) (sql.Result, error)
		FindOne(ctx context.Context, addressId uint64) (*OrderAddresses, error)
		Update(ctx context.Context, data *OrderAddresses) error
		Delete(ctx context.Context, addressId uint64) error
	}

	defaultOrderAddressesModel struct {
		sqlc.CachedConn
		table string
	}

	OrderAddresses struct {
		AddressId       uint64         `db:"address_id"`
		RecipientName   string         `db:"recipient_name"`   // 收件人姓名
		PhoneNumber     sql.NullString `db:"phone_number"`     // 联系电话
		Province        sql.NullString `db:"province"`         // 州/省
		City            string         `db:"city"`             // 城市
		DetailedAddress string         `db:"detailed_address"` // 详细地址
		CreatedAt       time.Time      `db:"created_at"`       // 创建时间
		UpdatedAt       time.Time      `db:"updated_at"`       // 更新时间
	}
)

func newOrderAddressesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultOrderAddressesModel {
	return &defaultOrderAddressesModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`order_addresses`",
	}
}

func (m *defaultOrderAddressesModel) Delete(ctx context.Context, addressId uint64) error {
	orderAddressesAddressIdKey := fmt.Sprintf("%s%v", cacheOrderAddressesAddressIdPrefix, addressId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `address_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, addressId)
	}, orderAddressesAddressIdKey)
	return err
}

func (m *defaultOrderAddressesModel) FindOne(ctx context.Context, addressId uint64) (*OrderAddresses, error) {
	orderAddressesAddressIdKey := fmt.Sprintf("%s%v", cacheOrderAddressesAddressIdPrefix, addressId)
	var resp OrderAddresses
	err := m.QueryRowCtx(ctx, &resp, orderAddressesAddressIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `address_id` = ? limit 1", orderAddressesRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, addressId)
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

func (m *defaultOrderAddressesModel) Insert(ctx context.Context, data *OrderAddresses) (sql.Result, error) {
	orderAddressesAddressIdKey := fmt.Sprintf("%s%v", cacheOrderAddressesAddressIdPrefix, data.AddressId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, orderAddressesRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.RecipientName, data.PhoneNumber, data.Province, data.City, data.DetailedAddress)
	}, orderAddressesAddressIdKey)
	return ret, err
}

func (m *defaultOrderAddressesModel) Update(ctx context.Context, data *OrderAddresses) error {
	orderAddressesAddressIdKey := fmt.Sprintf("%s%v", cacheOrderAddressesAddressIdPrefix, data.AddressId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `address_id` = ?", m.table, orderAddressesRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.RecipientName, data.PhoneNumber, data.Province, data.City, data.DetailedAddress, data.AddressId)
	}, orderAddressesAddressIdKey)
	return err
}

func (m *defaultOrderAddressesModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheOrderAddressesAddressIdPrefix, primary)
}

func (m *defaultOrderAddressesModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `address_id` = ? limit 1", orderAddressesRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultOrderAddressesModel) tableName() string {
	return m.table
}
