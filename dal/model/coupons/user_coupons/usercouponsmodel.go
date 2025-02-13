package user_coupons

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserCouponsModel = (*customUserCouponsModel)(nil)

type (
	// UserCouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCouponsModel.
	UserCouponsModel interface {
		userCouponsModel
		withSession(session sqlx.Session) UserCouponsModel
	}

	customUserCouponsModel struct {
		*defaultUserCouponsModel
	}
)

// NewUserCouponsModel returns a model for the database table.
func NewUserCouponsModel(conn sqlx.SqlConn) UserCouponsModel {
	return &customUserCouponsModel{
		defaultUserCouponsModel: newUserCouponsModel(conn),
	}
}

func (m *customUserCouponsModel) withSession(session sqlx.Session) UserCouponsModel {
	return NewUserCouponsModel(sqlx.NewSqlConnFromSession(session))
}
