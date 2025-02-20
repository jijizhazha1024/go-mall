package order_items

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OrderItemsModel = (*customOrderItemsModel)(nil)

type (
	// OrderItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemsModel.
	OrderItemsModel interface {
		orderItemsModel
		withSession(session sqlx.Session) OrderItemsModel
	}

	customOrderItemsModel struct {
		*defaultOrderItemsModel
	}
)

// NewOrderItemsModel returns a model for the database table.
func NewOrderItemsModel(conn sqlx.SqlConn) OrderItemsModel {
	return &customOrderItemsModel{
		defaultOrderItemsModel: newOrderItemsModel(conn),
	}
}

func (m *customOrderItemsModel) withSession(session sqlx.Session) OrderItemsModel {
	return NewOrderItemsModel(sqlx.NewSqlConnFromSession(session))
}
