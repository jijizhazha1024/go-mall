package coupon

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"time"
)

var _ CouponsModel = (*customCouponsModel)(nil)

type (
	// CouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponsModel.
	CouponsModel interface {
		couponsModel
		WithSession(session sqlx.Session) CouponsModel
		QueryCoupons(ctx context.Context, page, pageSize, ctype int32) ([]*Coupons, error)
		FindOneWithLock(ctx context.Context, session sqlx.Session, id string) (*Coupons, error)
		DecreaseStockWithSession(ctx context.Context, session sqlx.Session, id string, num int) error
		GetCouponTypeByID(ctx context.Context, session sqlx.Session, id string) (int64, error)
		CheckExpirationAndStatus(ctx context.Context, session sqlx.Session, id string) (bool, error)
	}

	customCouponsModel struct {
		*defaultCouponsModel
	}
)

func (m *customCouponsModel) CheckExpirationAndStatus(ctx context.Context, session sqlx.Session, id string) (bool, error) {
	var status CStatus
	query := fmt.Sprintf("select `status`, `end_time` from %s where id = ? limit 1", m.table)
	err := session.QueryRowCtx(ctx, &status, query, id)
	fmt.Println(status.Status)
	fmt.Println(status.EndTime)
	switch {
	case err == nil:
		return status.Status == 1 && status.EndTime.Before(time.Now()), nil
	case errors.Is(err, sqlx.ErrNotFound):
		return false, err
	default:
		return false, err
	}
}

func (m *customCouponsModel) GetCouponTypeByID(ctx context.Context, session sqlx.Session, id string) (int64, error) {
	var ctp int64
	query := fmt.Sprintf("select `type` from %s where id = ? limit 1", m.table)
	err := session.QueryRowCtx(ctx, &ctp, query, id)
	return ctp, err

}

func (m *customCouponsModel) FindOneWithLock(ctx context.Context, session sqlx.Session, id string) (*Coupons, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `id` = ? FOR SHARE", couponsRows, m.table)
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

func (m *customCouponsModel) WithSession(session sqlx.Session) CouponsModel {
	return NewCouponsModel(sqlx.NewSqlConnFromSession(session))
}
