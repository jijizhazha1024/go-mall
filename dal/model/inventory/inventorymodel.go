package inventory

import (
	"context"
	"errors"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InventoryModel = (*customInventoryModel)(nil)

type (
	// InventoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInventoryModel.
	InventoryModel interface {
		inventoryModel
		FindAll(ctx context.Context) ([]*Inventory, error)
		WithSession(session sqlx.Session) InventoryModel
		UpdateOrCreate(ctx context.Context, inventory Inventory) error
		DecreaseInventoryAtom(ctx context.Context, productId int32, quantity int32) (cnt int64, err error)
		ReturnInventory(ctx context.Context, id int32, quantity int32) (cnt int64, err error)
	}

	customInventoryModel struct {
		*defaultInventoryModel
	}
)

func (m *customInventoryModel) ReturnInventory(ctx context.Context, productId int32, quantity int32) (cnt int64, err error) {
	var inventory Inventory
	query := fmt.Sprintf("select * from %s where `product_id` = ? for update", m.table)
	if err := m.conn.QueryRowCtx(ctx, &inventory, query, productId); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return 0, err
		}
		return 0, biz.InventoryDecreaseFailedErr
	}
	cnt = inventory.Total + int64(quantity)
	query = fmt.Sprintf("UPDATE %s SET sold = sold - ?, total = total + ? WHERE product_id = ?", m.table)
	res, err := m.conn.ExecCtx(ctx, query, quantity, quantity, productId)
	if err != nil {
		return 0, biz.InventoryDecreaseFailedErr
	}
	if affected, err := res.RowsAffected(); err != nil {
		return 0, biz.InventoryDecreaseFailedErr
	} else if affected == 0 {
		return 0, biz.InventoryDecreaseFailedErr
	}
	return cnt, nil
}

func (m *customInventoryModel) UpdateOrCreate(ctx context.Context, inventory Inventory) error {
	var cnt int64
	query := fmt.Sprintf("select count(*) from %s where `product_id` = ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &cnt, query, inventory.ProductId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			_, err := m.Insert(ctx, &inventory)
			if err != nil {
				return err
			}
		}
		return err
	}

	return m.Update(ctx, &inventory)
}

func (m *customInventoryModel) DecreaseInventoryAtom(ctx context.Context, productId int32, quantity int32) (int64, error) {
	var cnt int64
	if err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {

		// --------------- check ---------------

		var inventory Inventory
		query := fmt.Sprintf("select * from %s where `product_id` = ? for update", m.table)
		if err := session.QueryRowCtx(ctx, &inventory, query, productId); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return err
			}
			return biz.InventoryDecreaseFailedErr
		}
		cnt = inventory.Total - inventory.Sold - int64(quantity)
		if cnt < int64(quantity) {
			return biz.InventoryNotEnoughErr
		}

		// --------------- update ---------------

		query = fmt.Sprintf("UPDATE %s SET sold = sold + ?, total = total - ? WHERE product_id = ?", m.table)
		res, err := session.ExecCtx(ctx, query, quantity, quantity, productId)
		if err != nil {
			return biz.InventoryDecreaseFailedErr
		}
		if affected, err := res.RowsAffected(); err != nil {
			return biz.InventoryDecreaseFailedErr
		} else if affected == 0 {
			return biz.InventoryDecreaseFailedErr
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return cnt, nil
}
func (m *customInventoryModel) FindAll(ctx context.Context) ([]*Inventory, error) {
	// 1. 构建 SQL 查询语
	var inventorys []*Inventory
	query := fmt.Sprintf("select * from %s ", m.table)
	err := m.conn.QueryRowsCtx(ctx, &inventorys, query)
	if err != nil {
		return nil, err
	}
	return inventorys, nil
}

// NewInventoryModel returns a model for the database table.
func NewInventoryModel(conn sqlx.SqlConn) InventoryModel {
	return &customInventoryModel{
		defaultInventoryModel: newInventoryModel(conn),
	}
}

func (m *customInventoryModel) WithSession(session sqlx.Session) InventoryModel {
	return NewInventoryModel(sqlx.NewSqlConnFromSession(session))
}
