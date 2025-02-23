package inventory

import (
	"context"
	"errors"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"strings"

	sqlx1 "github.com/jmoiron/sqlx"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InventoryModel = (*customInventoryModel)(nil)

type (
	// InventoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInventoryModel.
	InventoryModel interface {
		inventoryModel
		FindAll(ctx context.Context) ([]*Inventory, error)
		FindLockOrder(ctx context.Context, session sqlx.Session, order_id string, user_id int64, table string) (bool, error)
		LockOrder(ctx context.Context, session sqlx.Session, order_id string, user_id int64, table string) error
		WithSession(session sqlx.Session) InventoryModel
		UpdateOrCreate(ctx context.Context, inventory Inventory) error
		BatchReturn(ctx context.Context, session sqlx.Session, productIDs []int32, quantities []int32) error
		DecreaseInventoryAtom(ctx context.Context, productId int32, quantity int32) (cnt int64, err error)
		Batchdecrease(ctx context.Context, session sqlx.Session, productIDs []int32, quantities []int32) error
		BatchReturnInventoryAtom(ctx context.Context, productIDs []int32, quantities []int32, orderID string, userID int64) error
		ReturnInventory(ctx context.Context, id int32, quantity int32) (cnt int64, err error)
		BatchDecreaseInventoryAtom(ctx context.Context, productId []int32, quantity []int32, user_id int64, order_id string) error
	}

	customInventoryModel struct {
		*defaultInventoryModel
	}
)

func (m *customInventoryModel) BatchReturnInventoryAtom(ctx context.Context, productIDs []int32, quantities []int32, orderID string, userID int64) error {

	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {

		// 阶段1: 幂等检查
		isLocked, err := m.FindLockOrder(ctx, session, orderID, userID, m.lockreturntable)
		if err != nil {
			return err
		}
		if isLocked {

			return fmt.Errorf("订单 %s 已被锁定", orderID)
		}

		// 阶段2: 创建锁记录（30分钟有效期）
		if err := m.LockOrder(ctx, session, orderID, userID, m.lockreturntable); err != nil {
			return nil
		}

		// 阶段3: 批量锁定库存记录
		query := fmt.Sprintf(`
		SELECT product_id, total, sold 
		FROM %s 
		WHERE product_id IN (?)
		FOR UPDATE  -- 行级锁
	`, m.table)

		query, args, err := sqlx1.In(query, productIDs)
		if err != nil {
			return err
		}

		var inventories []*Inventory
		if err := session.QueryRowsCtx(ctx, &inventories, query, args...); err != nil {
			return err
		}
		//阶段4 ：批量更新库存
		err = m.BatchReturn(ctx, session, productIDs, quantities)
		if err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return err
	}
	return nil

}

func (m *customInventoryModel) BatchReturn(ctx context.Context, session sqlx.Session, productIDs []int32, quantities []int32) error {

	// 阶段3: 批量更新
	var updateBuilder strings.Builder
	updateBuilder.WriteString(fmt.Sprintf("UPDATE %s SET ", m.table))
	updateBuilder.WriteString("sold = CASE product_id ")

	// 构建WHEN条件
	whenCases := make([]string, len(productIDs))
	for i, pid := range productIDs {
		whenCases[i] = fmt.Sprintf("WHEN %d THEN sold - %d", pid, quantities[i])
	}
	updateBuilder.WriteString(strings.Join(whenCases, " "))
	updateBuilder.WriteString(" END, total = CASE product_id ")

	whenCases = whenCases[:0]
	for i, pid := range productIDs {
		whenCases = append(whenCases, fmt.Sprintf("WHEN %d THEN total + %d", pid, quantities[i]))
	}
	updateBuilder.WriteString(strings.Join(whenCases, " "))
	updateBuilder.WriteString(" END WHERE product_id IN (?)")

	// 执行更新
	updateQuery, updateArgs, err := sqlx1.In(updateBuilder.String(), productIDs)
	if err != nil {
		return err
	}

	res, err := session.ExecCtx(ctx, updateQuery, updateArgs...)
	if err != nil {
		return biz.InventoryDecreaseFailedErr
	}

	// 验证影响行数
	if affected, _ := res.RowsAffected(); affected != int64(len(productIDs)) {
		return biz.InventoryDecreaseFailedErr
	}

	return nil

}
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
	var exists bool
	query := fmt.Sprintf("select exists(select 1 from %s where `product_id` = ?)", m.table)
	err := m.conn.QueryRowCtx(ctx, &exists, query, inventory.ProductId)
	if err != nil {
		return err
	}

	if !exists {
		_, err := m.Insert(ctx, &inventory)
		if err != nil {
			return err
		}
		return nil
	}

	return m.Update(ctx, &inventory)
}

func (m *customInventoryModel) LockOrder(
	ctx context.Context,
	session sqlx.Session,
	orderID string,
	userID int64,
	table string,
) error {

	query := fmt.Sprintf("INSERT INTO  %s  (order_id, user_id) VALUES (?, ?)", table)
	_, err := session.ExecCtx(ctx, query, orderID, userID)
	if err != nil {
		return err
	}

	return nil
}

// FindLockOrder 幂等性检查
func (m *customInventoryModel) FindLockOrder(
	ctx context.Context,
	session sqlx.Session,
	orderID string,
	userID int64,
	table string,
) (bool, error) {
	// 构建 SQL 查询语
	query := fmt.Sprintf(`
	SELECT COUNT(*) 
	FROM %s 
	WHERE order_id = ? AND user_id = ? 
	FOR UPDATE`,
		table)

	var count int
	err := session.QueryRowCtx(ctx, &count, query, orderID, userID)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil

}
func (m *customInventoryModel) BatchDecreaseInventoryAtom(
	ctx context.Context,
	productIDs []int32,
	quantities []int32,
	userID int64,
	orderID string,
) error {

	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 阶段1: 幂等检查
		isLocked, err := m.FindLockOrder(ctx, session, orderID, userID, m.lockdecreasetable)
		if err != nil {
			return err
		}
		if isLocked {
			return nil
		}

		// 阶段2: 创建锁记录（30分钟有效期）
		if err := m.LockOrder(ctx, session, orderID, userID, m.lockdecreasetable); err != nil {
			return fmt.Errorf("创建锁失败: %w", err)
		}

		// --- 阶段3: 批量锁定库存记录 ---
		query := fmt.Sprintf(`
            SELECT product_id, total, sold 
            FROM %s 
            WHERE product_id IN (?)
            FOR UPDATE`, m.table) // 排他锁

		// 处理IN查询参数化
		query, args, err := sqlx1.In(query, productIDs)
		if err != nil {
			return fmt.Errorf("build IN query failed: %w", err)
		}

		var inventories []*Inventory
		if err := session.QueryRowsCtx(ctx, &inventories, query, args...); err != nil {
			return fmt.Errorf("batch lock inventory failed: %w", err)
		}

		// 转换为快速查找map
		inventoryMap := make(map[int32]*Inventory, len(inventories))
		for _, inv := range inventories {
			inventoryMap[int32(inv.ProductId)] = inv
		}

		// --- 阶段4: 库存预检查 ---
		for i, pid := range productIDs {
			inv, exists := inventoryMap[pid]
			if !exists {
				return fmt.Errorf("product %d not found: %w", pid, sqlx.ErrNotFound)
			}
			if inv.Total < int64(quantities[i]) {
				return fmt.Errorf("product %d not enough: %w", pid, biz.InventoryNotEnoughErr)
			}
		}

		// --- 阶段5: 执行批量扣减 ---
		if err := m.Batchdecrease(ctx, session, productIDs, quantities); err != nil {
			return fmt.Errorf("batch decrease failed: %w", err)
		}

		return nil
	})
}
func (m *customInventoryModel) Batchdecrease(ctx context.Context, session sqlx.Session, productIDs []int32, quantities []int32) error {
	// 阶段3: 批量更新
	var updateBuilder strings.Builder
	updateBuilder.WriteString(fmt.Sprintf("UPDATE %s SET ", m.table))
	updateBuilder.WriteString("sold = CASE product_id ")

	// 构建WHEN条件
	whenCases := make([]string, len(productIDs))
	for i, pid := range productIDs {
		whenCases[i] = fmt.Sprintf("WHEN %d THEN sold + %d", pid, quantities[i])
	}
	updateBuilder.WriteString(strings.Join(whenCases, " "))
	updateBuilder.WriteString(" END, total = CASE product_id ")

	whenCases = whenCases[:0]
	for i, pid := range productIDs {
		whenCases = append(whenCases, fmt.Sprintf("WHEN %d THEN total - %d", pid, quantities[i]))
	}
	updateBuilder.WriteString(strings.Join(whenCases, " "))
	updateBuilder.WriteString(" END WHERE product_id IN (?)")

	// 执行更新
	updateQuery, updateArgs, err := sqlx1.In(updateBuilder.String(), productIDs)
	if err != nil {
		return err
	}

	res, err := session.ExecCtx(ctx, updateQuery, updateArgs...)
	if err != nil {
		return biz.InventoryDecreaseFailedErr
	}

	// 验证影响行数
	if affected, _ := res.RowsAffected(); affected != int64(len(productIDs)) {
		return biz.InventoryDecreaseFailedErr
	}

	return nil
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
