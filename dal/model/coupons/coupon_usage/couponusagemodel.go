package coupon_usage

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CouponUsageModel = (*customCouponUsageModel)(nil)

type (
	// CouponUsageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponUsageModel.
	CouponUsageModel interface {
		couponUsageModel
		WithSession(session sqlx.Session) CouponUsageModel
		QueryUsageListByUserId(ctx context.Context, userId uint64, page, size int32) ([]*CouponUsage, error)
	}

	customCouponUsageModel struct {
		*defaultCouponUsageModel
	}
)

func (m *customCouponUsageModel) QueryUsageListByUserId(ctx context.Context, userId uint64, page, size int32) ([]*CouponUsage, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ? LIMIT ? OFFSET ?", couponUsageRows, m.table)
	resp := make([]*CouponUsage, 0)
	//pageSize, (page-1)*pageSize
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, size, (page-1)*size)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return resp, nil
	default:
		return nil, err
	}
}

// NewCouponUsageModel returns a model for the database table.
func NewCouponUsageModel(conn sqlx.SqlConn) CouponUsageModel {
	return &customCouponUsageModel{
		defaultCouponUsageModel: newCouponUsageModel(conn),
	}
}

func (m *customCouponUsageModel) WithSession(session sqlx.Session) CouponUsageModel {
	return NewCouponUsageModel(sqlx.NewSqlConnFromSession(session))
}
