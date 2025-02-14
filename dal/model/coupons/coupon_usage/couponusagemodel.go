package coupon_usage

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CouponUsageModel = (*customCouponUsageModel)(nil)

type (
	// CouponUsageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponUsageModel.
	CouponUsageModel interface {
		couponUsageModel
		WithSession(session sqlx.Session) CouponUsageModel
	}

	customCouponUsageModel struct {
		*defaultCouponUsageModel
	}
)

// NewCouponUsageModel returns a model for the database table.
func NewCouponUsageModel(conn sqlx.SqlConn) CouponUsageModel {
	return &customCouponUsageModel{
		defaultCouponUsageModel: newCouponUsageModel(conn),
	}
}

func (m *customCouponUsageModel) WithSession(session sqlx.Session) CouponUsageModel {
	return NewCouponUsageModel(sqlx.NewSqlConnFromSession(session))
}
