package order

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OrdersModel = (*customOrdersModel)(nil)

type (
	// OrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrdersModel.
	OrdersModel interface {
		ordersModel
		withSession(session sqlx.Session) OrdersModel
	}

	customOrdersModel struct {
		*defaultOrdersModel
	}
)

// NewOrdersModel returns a model for the database table.
func NewOrdersModel(conn sqlx.SqlConn) OrdersModel {
	return &customOrdersModel{
		defaultOrdersModel: newOrdersModel(conn),
	}
}

func (m *customOrdersModel) withSession(session sqlx.Session) OrdersModel {
	return NewOrdersModel(sqlx.NewSqlConnFromSession(session))
}
