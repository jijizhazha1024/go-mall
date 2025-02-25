package payment

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"strings"
)

var _ PaymentsModel = (*customPaymentsModel)(nil)

type (
	// PaymentsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentsModel.
	PaymentsModel interface {
		paymentsModel
		withSession(session sqlx.Session) PaymentsModel
		UpdateInfoByOrderId(ctx context.Context, newData *Payments) error
		Count(ctx context.Context) (int64, error)
		FindPage(ctx context.Context, userId uint32, offset, limit int) ([]*Payments, error)
		FindOneByOrderId(ctx context.Context, pre_order_id string) (*Payments, error)
	}

	customPaymentsModel struct {
		*defaultPaymentsModel
	}
)

// NewPaymentsModel returns a model for the database table.
func NewPaymentsModel(conn sqlx.SqlConn) PaymentsModel {
	return &customPaymentsModel{
		defaultPaymentsModel: newPaymentsModel(conn),
	}
}

var (
	paymentsRowsWithHolder = strings.Join(stringx.Remove(paymentsFieldNames, "`pre_order_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

func (m *customPaymentsModel) withSession(session sqlx.Session) PaymentsModel {
	return NewPaymentsModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultPaymentsModel) UpdateInfoByOrderId(ctx context.Context, newData *Payments) error {
	// 定义需要更新的字段
	paymentsRowsWithHolder := "`order_id`=?, `transaction_id`=?, `status`=?, `updated_at`=?, `paid_at`=?"

	// 构造 SQL 更新语句
	query := fmt.Sprintf("update %s set %s where `pre_order_id` = ?", m.table, paymentsRowsWithHolder)

	// 执行更新操作
	_, err := m.conn.ExecCtx(ctx, query,
		newData.OrderId.String,       // order_id
		newData.TransactionId.String, // transaction_id
		newData.Status,               // status
		newData.UpdatedAt,            // updated_at
		newData.PaidAt.Int64,         // paid_at
		newData.PreOrderId,           // WHERE pre_order_id
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
func (m *defaultPaymentsModel) FindOneByOrderId(ctx context.Context, pre_order_id string) (*Payments, error) {
	query := fmt.Sprintf("select %s from %s where `pre_order_id` = ? limit 1", paymentsRows, m.table)
	var resp Payments
	err := m.conn.QueryRowCtx(ctx, &resp, query, pre_order_id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
