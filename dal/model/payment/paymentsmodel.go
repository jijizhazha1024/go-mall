package payment

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentsModel = (*customPaymentsModel)(nil)

type (
	// PaymentsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentsModel.
	PaymentsModel interface {
		paymentsModel
		WithSession(session sqlx.Session) PaymentsModel
		UpdateInfoByOrderId(ctx context.Context, newData *Payments) error
		Count(ctx context.Context) (int64, error)
		FindPage(ctx context.Context, userId uint32, offset, limit int) ([]*Payments, error)
		FindOneByOrderId(ctx context.Context, pre_order_id string) (*Payments, error)
		CheckExistByOrderId(ctx context.Context, orderID string) (bool, error)
	}

	customPaymentsModel struct {
		*defaultPaymentsModel
	}
)

func (m *customPaymentsModel) CheckExistByOrderId(ctx context.Context, orderID string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE `pre_order_id` = ?", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, orderID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// NewPaymentsModel returns a model for the database table.
func NewPaymentsModel(conn sqlx.SqlConn) PaymentsModel {
	return &customPaymentsModel{
		defaultPaymentsModel: newPaymentsModel(conn),
	}
}

func (m *customPaymentsModel) WithSession(session sqlx.Session) PaymentsModel {
	return NewPaymentsModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultPaymentsModel) UpdateInfoByOrderId(ctx context.Context, newData *Payments) error {
	// 定义需要更新的字段
	paymentsRowsWithHolder := "`transaction_id`=?, `status`=?, `paid_at`=?"

	// 构造 SQL 更新语句
	query := fmt.Sprintf("update %s set %s where `order_id` = ?", m.table, paymentsRowsWithHolder)

	// 执行更新操作
	_, err := m.conn.ExecCtx(ctx, query,
		newData.TransactionId,
		newData.Status,
		newData.PaidAt,
		newData.OrderId,
	)
	return err
}

// 查询支付记录
func (m *defaultPaymentsModel) FindPage(ctx context.Context, userId uint32, offset, limit int) ([]*Payments, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE `user_id` = ? LIMIT ? OFFSET ?", m.table)
	var payments []*Payments
	err := m.conn.QueryRowsCtx(ctx, &payments, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return payments, nil
}
func (m *defaultPaymentsModel) Count(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (m *defaultPaymentsModel) FindOneByOrderId(ctx context.Context, orderID string) (*Payments, error) {
	query := fmt.Sprintf("select %s from %s where `order_id` = ? limit 1 for share", paymentsRows, m.table)
	var resp Payments
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderID)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
