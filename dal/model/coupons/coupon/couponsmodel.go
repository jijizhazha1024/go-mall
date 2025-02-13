package coupon

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CouponsModel = (*customCouponsModel)(nil)

type (
	// CouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponsModel.
	CouponsModel interface {
		couponsModel
		withSession(session sqlx.Session) CouponsModel
	}

	customCouponsModel struct {
		*defaultCouponsModel
	}
)

// NewCouponsModel returns a model for the database table.
func NewCouponsModel(conn sqlx.SqlConn) CouponsModel {
	return &customCouponsModel{
		defaultCouponsModel: newCouponsModel(conn),
	}
}

func (m *customCouponsModel) withSession(session sqlx.Session) CouponsModel {
	return NewCouponsModel(sqlx.NewSqlConnFromSession(session))
}
