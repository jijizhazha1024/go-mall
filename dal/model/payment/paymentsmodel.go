package payment

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PaymentsModel = (*customPaymentsModel)(nil)

type (
	// PaymentsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentsModel.
	PaymentsModel interface {
		paymentsModel
		withSession(session sqlx.Session) PaymentsModel
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

func (m *customPaymentsModel) withSession(session sqlx.Session) PaymentsModel {
	return NewPaymentsModel(sqlx.NewSqlConnFromSession(session))
}
