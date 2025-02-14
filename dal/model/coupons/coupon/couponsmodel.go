package coupon

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ CouponsModel = (*customCouponsModel)(nil)

type (
	// CouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponsModel.
	CouponsModel interface {
		couponsModel
		withSession(session sqlx.Session) CouponsModel
		QueryCoupons(ctx context.Context, page, pageSize, ctype int32) ([]*Coupons, error)
		FindOneWithLock(ctx context.Context, session sqlx.Session, id string) (*Coupons, error)
		DecreaseStockWithSession(ctx context.Context, session sqlx.Session, id string, num int) error
	}

	customCouponsModel struct {
		*defaultCouponsModel
	}
)

func (m *customCouponsModel) FindOneWithLock(ctx context.Context, session sqlx.Session, id string) (*Coupons, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `id` = ? FOR UPDATE", couponsRows, m.table)
	var resp Coupons
	err := session.QueryRowCtx(ctx, &resp, query, id)
	return &resp, err
}

func (m *customCouponsModel) DecreaseStockWithSession(ctx context.Context, session sqlx.Session, id string, num int) error {
	query := fmt.Sprintf(
		"UPDATE %s SET remaining_count = remaining_count - ? WHERE id = ? AND remaining_count >= ?",
		m.table,
	)
	_, err := session.ExecCtx(ctx, query, num, id, num)
	return err
}

func (m *customCouponsModel) QueryCoupons(ctx context.Context, page, pageSize, ctype int32) ([]*Coupons, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", couponsRows, m.table)

	// 构建WHERE条件
	var where []string
	var args []interface{}

	if ctype != 0 {
		where = append(where, "type = ?")
		args = append(args, ctype)
	}

	// 组合WHERE条件
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// 添加分页
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	var coupons []*Coupons
	err := m.conn.QueryRowsCtx(ctx, &coupons, query, args...)
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

// NewCouponsModel returns a model for the database table.
func NewCouponsModel(conn sqlx.SqlConn) CouponsModel {
	return &customCouponsModel{
		defaultCouponsModel: newCouponsModel(conn),
	}
}

func (m *customCouponsModel) withSession(session sqlx.Session) CouponsModel {
	return NewCouponsModel(sqlx.NewSqlConnFromSession(session))
}
