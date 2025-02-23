// Code generated by goctl. DO NOT EDIT.

package inventory

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
	inventoryFieldNames          = builder.RawFieldNames(&Inventory{})
	inventoryRows                = strings.Join(inventoryFieldNames, ",")
	inventoryRowsExpectAutoSet   = strings.Join(stringx.Remove(inventoryFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	inventoryRowsWithPlaceHolder = strings.Join(stringx.Remove(inventoryFieldNames, "`product_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)



type (
	inventoryModel interface {
		Insert(ctx context.Context, data *Inventory) (sql.Result, error)
		FindOne(ctx context.Context, productId int64) (*Inventory, error)
		Update(ctx context.Context, data *Inventory) error
		Delete(ctx context.Context, productId int64) error
	}

	defaultInventoryModel struct {
		conn  sqlx.SqlConn
		table string
		lockdecreasetable string
		lockreturntable string
	}

	Inventory struct {
		ProductId int64 `db:"product_id"`
		Total     int64 `db:"total"`
		Sold      int64 `db:"sold"`
	}
)

func newInventoryModel(conn sqlx.SqlConn) *defaultInventoryModel {
	return &defaultInventoryModel{
		conn:  conn,
		table: "`inventory`",
		lockdecreasetable:      "inventory_lock",  
		lockreturntable:       "return_lock",  
	}
}

func (m *defaultInventoryModel) Delete(ctx context.Context, productId int64) error {
	query := fmt.Sprintf("delete from %s where `product_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, productId)
	return err
}

func (m *defaultInventoryModel) FindOne(ctx context.Context, productId int64) (*Inventory, error) {
	
	query := fmt.Sprintf("select %s from %s where `product_id` = ? limit 1", inventoryRows, m.table)
	var resp Inventory
	err := m.conn.QueryRowCtx(ctx, &resp, query, productId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultInventoryModel) Insert(ctx context.Context, data *Inventory) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, inventoryRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ProductId, data.Total, data.Sold)
	return ret, err
}

func (m *defaultInventoryModel) Update(ctx context.Context, data *Inventory) error {
	query := fmt.Sprintf("update %s set %s where `product_id` = ?", m.table, inventoryRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Total, data.Sold, data.ProductId)
	return err
}

func (m *defaultInventoryModel) tableName() string {
	return m.table
}
