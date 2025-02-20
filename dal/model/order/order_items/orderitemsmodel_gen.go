// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.5

package order_items

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
	orderItemsFieldNames          = builder.RawFieldNames(&OrderItems{})
	orderItemsRows                = strings.Join(orderItemsFieldNames, ",")
	orderItemsRowsExpectAutoSet   = strings.Join(stringx.Remove(orderItemsFieldNames, "`item_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	orderItemsRowsWithPlaceHolder = strings.Join(stringx.Remove(orderItemsFieldNames, "`item_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	orderItemsModel interface {
		Insert(ctx context.Context, data *OrderItems) (sql.Result, error)
		FindOne(ctx context.Context, itemId uint64) (*OrderItems, error)
		Update(ctx context.Context, data *OrderItems) error
		Delete(ctx context.Context, itemId uint64) error
	}

	defaultOrderItemsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	OrderItems struct {
		ItemId      uint64         `db:"item_id"`
		OrderId     string         `db:"order_id"`     // 关联订单号
		ProductId   int64          `db:"product_id"`   // 商品ID
		Quantity    int64          `db:"quantity"`     // 购买数量
		ProductName string         `db:"product_name"` // 商品名称
		ProductDesc sql.NullString `db:"product_desc"` // 规格描述
		UnitPrice   int64          `db:"unit_price"`   // 单价(分)
	}
)

func newOrderItemsModel(conn sqlx.SqlConn) *defaultOrderItemsModel {
	return &defaultOrderItemsModel{
		conn:  conn,
		table: "`order_items`",
	}
}

func (m *defaultOrderItemsModel) Delete(ctx context.Context, itemId uint64) error {
	query := fmt.Sprintf("delete from %s where `item_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, itemId)
	return err
}

func (m *defaultOrderItemsModel) FindOne(ctx context.Context, itemId uint64) (*OrderItems, error) {
	query := fmt.Sprintf("select %s from %s where `item_id` = ? limit 1", orderItemsRows, m.table)
	var resp OrderItems
	err := m.conn.QueryRowCtx(ctx, &resp, query, itemId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultOrderItemsModel) Insert(ctx context.Context, data *OrderItems) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, orderItemsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.OrderId, data.ProductId, data.Quantity, data.ProductName, data.ProductDesc, data.UnitPrice)
	return ret, err
}

func (m *defaultOrderItemsModel) Update(ctx context.Context, data *OrderItems) error {
	query := fmt.Sprintf("update %s set %s where `item_id` = ?", m.table, orderItemsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.OrderId, data.ProductId, data.Quantity, data.ProductName, data.ProductDesc, data.UnitPrice, data.ItemId)
	return err
}

func (m *defaultOrderItemsModel) tableName() string {
	return m.table
}
