package inventory_lock

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ InventoryLockModel = (*customInventoryLockModel)(nil)

type (
	// InventoryLockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInventoryLockModel.
	InventoryLockModel interface {
		inventoryLockModel
		withSession(session sqlx.Session) InventoryLockModel
	}

	customInventoryLockModel struct {
		*defaultInventoryLockModel
	}
)

// NewInventoryLockModel returns a model for the database table.
func NewInventoryLockModel(conn sqlx.SqlConn) InventoryLockModel {
	return &customInventoryLockModel{
		defaultInventoryLockModel: newInventoryLockModel(conn),
	}
}

func (m *customInventoryLockModel) withSession(session sqlx.Session) InventoryLockModel {
	return NewInventoryLockModel(sqlx.NewSqlConnFromSession(session))
}
